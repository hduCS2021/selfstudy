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
	left, err := IsTodayLeaveByID(stu.StudentId)
	if err != nil {
		return err
	}
	if left {
		return errors.New("已经请假过了")
	}
	leave := &TempLeave{StudentId: getStuByName(stu.Name).StudentId, Reason: reason, LeaveDay: leaveDay}
	if errs := db.Create(leave).GetErrors(); len(errs) != 0 {
		return errs2err(errs)
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

	if errs := db.Create(&leave).GetErrors(); len(errs) != 0 {
		return errs2err(errs)
	}
	return nil
}
