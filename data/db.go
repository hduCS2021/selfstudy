package data

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	if err := db.First(&stu, l.StudentId).Error; err != nil {
		return nil, err
	}
	return &stu, nil
}

func (l TempLeave) GetReason() string {
	return l.Reason
}

type LongLeave struct {
	Student   Student `gorm:"foreignKey:StudentId"`
	StudentId int32   `gorm:"uniqueIndex:one_leave_one_day"`
	Reason    string
	WeekDay   int32 `gorm:"uniqueIndex:one_leave_one_day"`
	Single    int32 `gorm:"uniqueIndex:one_leave_one_day"`
	gorm.Model
}

func (l LongLeave) GetStu() (*Student, error) {
	stu := Student{}
	if err := db.First(&stu, l.StudentId).Error; err != nil {
		return nil, err
	}
	return &stu, nil
}

func (l LongLeave) GetReason() string {
	return l.Reason
}

type Student struct {
	Name       string
	StudentId  int32 `gorm:"primaryKey"`
	QQ         int64
	Permission int32
}

type Checkin struct {
	Student   Student `gorm:"foreignKey:StudentId"`
	StudentId int32
	CheckinAt time.Time
}

func InitDB(source string) error {
	database, err := gorm.Open(mysql.Open(source), &gorm.Config{})
	if err != nil {
		return err
	}
	db = database
	if !db.Migrator().HasTable(&Student{}) {
		if err := database.Migrator().CreateTable(&Student{}).Error(); err != "" {
			return fmt.Errorf("can't create database:%v", err)
		}
	}
	if !db.Migrator().HasTable(&TempLeave{}) {
		if err := database.Migrator().CreateTable(&TempLeave{}).Error(); err != "" {
			return fmt.Errorf("can't create database:%v", err)
		}
	}
	if !db.Migrator().HasTable(&LongLeave{}) {
		if err := database.Migrator().CreateTable(&LongLeave{}).Error(); err != "" {
			return fmt.Errorf("can't create database:%v", err)
		}
	}
	if !db.Migrator().HasTable(&Checkin{}) {
		if err := database.Migrator().CreateTable(&Checkin{}).Error(); err != "" {
			return fmt.Errorf("can't create database:%v", err)
		}
	}
	return nil
}
