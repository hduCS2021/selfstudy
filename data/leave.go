package data

import (
	"fmt"
	"time"
)

func AddLeaveByQQ(qq int64, reason string, leaveDay time.Time) error {
	stu := getStuByQQ(qq)
	if stu.StudentId == 0 {
		return fmt.Errorf("can't find student")
	}
	leave := &Leave{StudentId: getStuByName(stu.Name).StudentId, Reason: reason, LeaveDay: leaveDay}
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
		WeekDay:   weekday,
		Single:    single,
	}

	if errs := db.Create(&leave).GetErrors(); len(errs) != 0 {
		return errs2err(errs)
	}
	return nil
}
