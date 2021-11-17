package data

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func errs2err(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return translateErr(errs[0])
	}
	var err = errs[0]
	for _, v := range errs[1:] {
		err = fmt.Errorf("%v;%v", err, translateErr(v))
	}
	return err
}

func translateErr(err error) error {
	if strings.Contains(err.Error(), "Duplicate entry") {
		return errors.New("记录重复")
	}
	return err
}

func sameDay(t1, t2 time.Time) bool {
	if t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year() {
		return true
	}
	return false
}
