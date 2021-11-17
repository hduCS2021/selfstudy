package data

import (
	"github.com/jinzhu/gorm"
	"time"
)

var TermBegin = time.Date(2021, 9, 13, 0, 0, 0, 0, time.Now().Location())

func QueryLeaveByID(id int32) ([]TempLeave, error) {
	var leaves []TempLeave
	errs := db.Where(&TempLeave{StudentId: id}).Find(&leaves).GetErrors()
	return leaves, errs2err(errs)
}

func QueryTodayLeaves() ([]Leave, error) {
	var tempLeaves []TempLeave
	var longLeaves []LongLeave
	var leaves []Leave
	n := time.Now()
	today := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	if errs := db.Where("leave_day = ?", today).Find(&tempLeaves).GetErrors(); len(errs) != 0 {
		return nil, errs2err(errs)
	}
	weeks := time.Since(TermBegin)/(time.Hour*24*7) + 1
	Ignore := (weeks+1)%2 + 1
	if errs := db.Where("week_day = ? AND single <> ?", n.Weekday(), Ignore).Find(&longLeaves).GetErrors(); len(errs) != 0 {
		return nil, errs2err(errs)
	}
	for _, v := range tempLeaves {
		leaves = append(leaves, v)
	}
	for _, v := range longLeaves {
		leaves = append(leaves, v)
	}
	return leaves, nil
}

func IsTodayLeaveByID(id int32) (bool, error) {
	var tempLeaves TempLeave
	var longLeaves LongLeave
	n := time.Now()
	today := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	if errs := db.Where("leave_day = ? AND student_id = ?", today, id).First(&tempLeaves).GetErrors(); len(errs) != 0 && !gorm.IsRecordNotFoundError(errs[0]) {
		return false, errs2err(errs)
	}
	weeks := time.Since(TermBegin)/(time.Hour*24*7) + 1
	Ignore := (weeks+1)%2 + 1
	if errs := db.Where("week_day = ? AND single <> ? AND student_id = ?", n.Weekday(), Ignore, id).First(&longLeaves).GetErrors(); len(errs) != 0 && !gorm.IsRecordNotFoundError(errs[0]) {
		return false, errs2err(errs)
	}
	if tempLeaves.StudentId != 0 || longLeaves.StudentId != 0 {
		return true, nil
	}
	return false, nil
}
