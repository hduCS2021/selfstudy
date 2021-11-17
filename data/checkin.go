package data

import (
	"errors"
	"time"
)

var ErrRepeat = errors.New("重复打卡")

func AddCheckinByID(id int32) error {
	t := time.Now()
	c := Checkin{}
	db.Where("CheckinAt = ?", t).First(&c)
	if c.Student.StudentId == id && sameDay(t, c.CheckinAt) {
		return ErrRepeat
	}
	return errs2err(db.Create(&Checkin{
		StudentId: id,
		CheckinAt: t,
	}).GetErrors())
}
