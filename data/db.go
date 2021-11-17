package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

type Leave interface {
	GetStu() (*Student, error)
	GetReason() string
}

type TempLeave struct {
	Student   Student `gorm:"foreignKey:StudentId"`
	StudentId int32   `gorm:"unique_index:one_leave_one_day"`
	Reason    string
	LeaveDay  time.Time `gorm:"unique_index:one_leave_one_day"`
	gorm.Model
}

func (l TempLeave) GetStu() (*Student, error) {
	stu := Student{}
	if errs := db.First(&stu, l.StudentId).GetErrors(); len(errs) != 0 {
		return nil, errs2err(errs)
	}
	return &stu, nil
}

func (l TempLeave) GetReason() string {
	return l.Reason
}

type LongLeave struct {
	Student   Student `gorm:"foreignKey:StudentId"`
	StudentId int32   `gorm:"unique_index:one_leave_one_day"`
	Reason    string
	WeekDay   int32 `gorm:"unique_index:one_leave_one_day"`
	Single    int32 `gorm:"unique_index:one_leave_one_day"`
	gorm.Model
}

func (l LongLeave) GetStu() (*Student, error) {
	stu := Student{}
	if err := db.First(&stu, l.StudentId).GetErrors(); len(err) != 0 {
		return nil, errs2err(err)
	}
	return &stu, nil
}

func (l LongLeave) GetReason() string {
	return l.Reason
}

type Student struct {
	Name       string
	StudentId  int32 `gorm:"primary_key"`
	QQ         int64
	Permission int32
}

type Checkin struct {
	Student   Student `gorm:"foreignKey:StudentId"`
	StudentId int32
	CheckinAt time.Time
}

func InitDB(source string) error {
	database, err := gorm.Open("mysql", source)
	if err != nil {
		return err
	}
	db = database
	if !db.HasTable(&Student{}) {
		if errs := database.CreateTable(&Student{}).GetErrors(); len(errs) != 0 {
			return fmt.Errorf("can't create database:%v", errs2err(errs))
		}
	}
	if !db.HasTable(&TempLeave{}) {
		if errs := database.CreateTable(&TempLeave{}).GetErrors(); len(errs) != 0 {
			return fmt.Errorf("can't create database:%v", errs2err(errs))
		}
	}
	if !db.HasTable(&LongLeave{}) {
		if errs := database.CreateTable(&LongLeave{}).GetErrors(); len(errs) != 0 {
			return fmt.Errorf("can't create database:%v", errs2err(errs))
		}
	}
	if !db.HasTable(&Checkin{}) {
		if errs := database.CreateTable(&Checkin{}).GetErrors(); len(errs) != 0 {
			return fmt.Errorf("can't create database:%v", errs2err(errs))
		}
	}
	return nil
}
