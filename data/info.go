package data

import "log"

func getStuByQQ(qq int64) *Student {
	stu := Student{}
	if err := db.First(&stu, "qq = ?", qq).Error; err != nil {
		log.Fatalln(err)
	}
	return &stu
}

func getStuByName(name string) *Student {
	stu := Student{}
	if err := db.First(&stu, "name = ?", name).Error; err != nil {
		log.Fatalln(err)
	}
	return &stu
}

func CheckQQ(qq int64) bool {
	if getStuByQQ(qq).StudentId == 0 {
		return false
	}
	return true
}
