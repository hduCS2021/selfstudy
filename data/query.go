package data

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

var TermBegin = time.Date(2021, 9, 13, 0, 0, 0, 0, time.Now().Location())

func QueryLeaveByID(id int32) ([]Leave, error) {
	var tempLeaves []TempLeave
	var longLeaves []LongLeave
	var leaves []Leave

	if err := db.Where(&TempLeave{StudentId: id}).Find(&tempLeaves).Error; err != nil {
		return nil, err
	}
	if err := db.Where(&TempLeave{StudentId: id}).Find(&longLeaves).Error; err != nil {
		return nil, err
	}

	for _, v := range tempLeaves {
		leaves = append(leaves, v)
	}
	for _, v := range longLeaves {
		leaves = append(leaves, v)
	}
	return leaves, nil
}

func QueryTodayLeaves() ([]Leave, error) {
	var tempLeaves []TempLeave
	var longLeaves []LongLeave
	var leaves []Leave
	n := time.Now()
	today := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	if err := db.Where("leave_day = ?", today).Find(&tempLeaves).Error; err != nil {
		return nil, err
	}
	weeks := int(time.Since(TermBegin)/(time.Hour*24*7) + 1)
	Ignore := 2 - (weeks+1)%2
	if err := db.Where("week_day = ? AND single <> ?", n.Weekday(), Ignore).Find(&longLeaves).Error; err != nil {
		return nil, err
	}
	for _, v := range tempLeaves {
		leaves = append(leaves, v)
	}
	for _, v := range longLeaves {
		leaves = append(leaves, v)
	}
	return leaves, nil
}

func IsTodayTempLeaveByID(id int32) (uint, error) {
	var tempLeaves TempLeave
	n := time.Now()
	today := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	if err := db.Where("leave_day = ? AND student_id = ?", today, id).First(&tempLeaves).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return tempLeaves.ID, nil

}

func IsTodayLongLeaveByID(id int32) (uint, error) {
	var longLeaves LongLeave
	n := time.Now()
	weeks := int(time.Since(TermBegin)/(time.Hour*24*7) + 1)
	Ignore := 2 - (weeks+1)%2
	if err := db.Where("week_day = ? AND single <> ? AND student_id = ?", n.Weekday(), Ignore, id).First(&longLeaves).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return longLeaves.ID, nil
}

func QueryShouldArrive() ([]Student, error) {
	n := time.Now()
	today := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	weeks := int(time.Since(TermBegin)/(time.Hour*24*7) + 1)
	Ignore := 2 - (weeks+1)%2
	var list []Student
	err := db.Where(
		"student_id NOT IN (? UNION ?)",
		db.Table("temp_leaves").
			Where("leave_day = ? ", today).
			Select("student_id"),
		db.Table("long_leaves").
			Where("week_day = ? AND single <> ?", int(n.Weekday()), Ignore).
			Select("student_id")).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
