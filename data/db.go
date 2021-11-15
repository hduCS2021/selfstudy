package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

type Leave struct {
	Student   Student `gorm:"foreignKey:StudentId"`
	StudentId int32
	Reason    string
	LeaveDay  time.Time
	gorm.Model
}

type LongLeave struct {
	Student   Student `gorm:"foreignKey:StudentId"`
	StudentId int32
	Reason    string
	WeekDay   int
	Single    int
	gorm.Model
}

type Student struct {
	Name      string
	StudentId int32 `gorm:"primaryKey"`
	QQ        int64
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
	if !db.HasTable(&Leave{}) {
		if errs := database.CreateTable(&Leave{}).GetErrors(); len(errs) != 0 {
			return fmt.Errorf("can't create database:%v", errs2err(errs))
		}
	}
	if !db.HasTable(&LongLeave{}) {
		if errs := database.CreateTable(&LongLeave{}).GetErrors(); len(errs) != 0 {
			return fmt.Errorf("can't create database:%v", errs2err(errs))
		}
	}
	return nil
}
