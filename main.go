package main

import (
	"bytes"
	"fmt"
	"github.com/BaiMeow/SimpleBot/bot"
	"github.com/BaiMeow/SimpleBot/driver"
	"github.com/BaiMeow/SimpleBot/message"
	"github.com/hduCS2021/selfstudy/data"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"time"
)

var DailyLeavesTmpl = template.Must(template.New("DailyLeaves").Parse(`今日请假名单{{ range $name,$reason :=  . }}
{{ $name }}:{{ $reason }}
{{ end }}`))

var b *bot.Bot

const help = `夜自修自助请假功能介绍：
①临时请假，仅可请假当天 “夜自修请假 <reason>”
②长期请假 “夜自修长期请假 <周几，可填'周一'，'周二'...’周日‘> <请假原因> <可填单周/双周，不填表示都请假>”
③调取班级今日请假记录 “今日请假” 
`

var week = []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
var singleWeek = []string{"单双周", "单周", "双周"}

func main() {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("toml")
	vp.AddConfigPath("./")
	vp.SetDefault("addr", "ws://localhost:7000")
	vp.SetDefault("passwd", "")
	vp.SetDefault("dbSource", "user:passwd@tcp(localhost:3306)/table")
	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := vp.WriteConfigAs("./config.toml")
			if err != nil {
				return
			}
			log.Fatalln("config not found")
		}
		log.Fatalln(err)
	}
	if err := data.InitDB(vp.GetString("dbSource")); err != nil {
		log.Fatal(err)
	}

	addr := vp.GetString("addr")
	passwd := vp.GetString("passwd")

	b = bot.New(driver.NewWsDriver(addr, passwd))
	b.Attach(&bot.PrivateMsgHandler{
		F:        handlePrivateMsg,
		Priority: 1,
	})
	b.Attach(&bot.FriendAddHandler{
		F:        AutoAccept,
		Priority: 1,
	})
	if err := b.Run(); err != nil {
		log.Fatal(err)
	}
	select {}
}

func handlePrivateMsg(_ int32, UserID int64, Msg message.Msg) bool {
	rec := func(txt string) {
		if _, err := b.SendPrivateMsg(UserID, message.New().Text(txt)); err != nil {
			log.Println(err)
		}
	}
	if !data.CheckQQ(UserID) {
		return false
	}
	msgs := Msg.Fields()
	if msgs[0].GetType() != "text" {
		return false
	}
	switch msgs[0].(message.Text).Text {
	case "夜自修请假":
		if len(msgs) != 2 && msgs[1].GetType() != "text" {
			rec(cmdError("缺少请假原因"))
			return true
		}
		reason := msgs[1].(message.Text).Text
		now := time.Now()
		if err := data.AddTempLeaveByQQ(UserID, reason, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())); err != nil {
			rec(fmt.Sprintf("数据库错误：%v", err))
			return true
		}
		t := now.Format("2006年1月2日")
		rec(fmt.Sprintf("请假成功:\n时间:%s\n原因:%s", t, reason))
	case "夜自修长期请假":
		if len(msgs) < 2 || msgs[1].GetType() != "text" {
			rec(cmdError("缺少请假时间"))
			return true
		}
		weekday := msgs[1].(message.Text).Text
		weekdayN := 0
		for i, v := range week {
			if weekday == v {
				weekdayN = i + 1
				break
			}
		}
		if weekdayN == 0 {
			rec(cmdError("请假时间无效"))
			return true
		}
		if len(msgs) < 3 || msgs[2].GetType() != "text" {
			rec(cmdError("缺少请假原因"))
			return true
		}
		reason := msgs[2].(message.Text).Text
		var single = 0
		if len(msgs) >= 4 && msgs[3].GetType() == "text" {
			kind := msgs[3].(message.Text).Text
			switch kind {
			case "单周":
				single = 1
			case "双周":
				single = 2
			default:
				rec(cmdError("要不填单周要不填双周要不就不填，不要填奇奇怪怪的东西上来"))
				return true
			}
		}
		if err := data.AddLongLeaveByQQ(UserID, reason, weekdayN, single); err != nil {
			rec(fmt.Sprintf("数据库错误：%v", err))
			return false
		}
		rec(fmt.Sprintf("请假成功：\n时间：%s %s\n原因：%s", week[weekdayN-1], singleWeek[single], reason))
	case "今日请假":
		if len(msgs) != 1 {
			rec(cmdError("如果想查询班级今日请假信息不要输入额外的东西"))
			return true
		}
		leaves, err := data.QueryTodayLeaves()
		if err != nil {
			rec(fmt.Sprintf("数据库错误：%v", err))
			log.Println(err)
			return true
		}
		if len(leaves) == 0 {
			rec("今日暂无请假记录")
			return true
		}

		list := make(map[string]string)
		for _, v := range leaves {
			stu, err := v.GetStu()
			if err != nil {
				rec(fmt.Sprintf("数据库错误：%v", err))
				log.Println(err)
			}

			if _, ok := list[stu.Name]; ok {
				//优先展示长时间请假
				if _, ok := v.(data.LongLeave); ok {
					list[stu.Name] = v.GetReason()
				}
				continue
			}
			list[stu.Name] = v.GetReason()

		}
		buffer := new(bytes.Buffer)
		if err := DailyLeavesTmpl.Execute(buffer, list); err != nil {
			rec(fmt.Sprintf("模板解析错误：%v", err))
			log.Println(err)
		}
		rec(buffer.String())

	case "帮助":
		if len(msgs) != 1 {
			rec(cmdError("如果想查询帮助信息不要输入额外的东西"))
			return true
		}
		rec(help)
		return true

	}
	return true
}

func cmdError(reason string) string {
	return fmt.Sprintf("命令错误：%s", reason)
}

func AutoAccept(request *bot.FriendRequest) bool {
	if !data.CheckQQ(request.UserID) {
		return false
	}
	request.Agree("")
	if _, err := b.SendPrivateMsg(request.UserID, message.New().Text(help)); err != nil {
		log.Println(err)
		return true
	}
	return true
}
