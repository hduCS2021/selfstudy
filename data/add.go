package data

import (
	"errors"
	"fmt"
	"time"
)

func AddTempLeaveByQQ(qq int64, reason string, leaveDay time.Time) error {
	stu := getStuByQQ(qq)
	if stu.StudentId == 0 {
		return fmt.Errorf("can't find student")
	}

	if id, err := IsTodayLongLeaveByID(stu.StudentId); err != nil {
		return err
	} else if id != 0 {
		return errors.New("已存在长期请假")
	}

	leave := &TempLeave{
		StudentId: getStuByName(stu.Name).StudentId,
		Reason:    reason,
		LeaveDay:  leaveDay,
	}

	if id, err := IsTodayTempLeaveByID(stu.StudentId); err != nil {
		return err
	} else if id != 0 {
		if err := db.Where("id = ?", id).Update("reason", reason).Error; err != nil {
			return err
		}
		return nil
	}
	if err := db.Create(leave).Error; err != nil {
		return err
	}
	return nil
}

func AddLongLeaveByQQ(qq int64, reason string, weekday int, single int) error {
	stu := getStuByQQ(qq)
	if stu.StudentId == 0 {
		return fmt.Errorf("can't find student")
	}
	leave := LongLeave{
		StudentId: stu.StudentId,
		Reason:    reason,
		WeekDay:   int32(weekday),
		Single:    int32(single),
	}
	if id, err := IsTodayLongLeaveByID(stu.StudentId); err != nil {
		return err
	} else if id != 0 {
		if err := db.Where("id = ?", id).Update("reason", reason).Error; err != nil {
			return err
		}
		return nil
	}
	if err := db.Create(&leave).Error; err != nil {
		return err
	}
	return nil
}
