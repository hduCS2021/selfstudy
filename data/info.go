package data

import "log"

func getStuByQQ(qq int64) *Student {
	stu := Student{}
	if errs := db.First(&stu, "qq = ?", qq).GetErrors(); len(errs) != 0 {
		log.Fatalln(errs)
	}
	return &stu
}

func getStuByName(name string) *Student {
	stu := Student{}
	if errs := db.First(&stu, "name = ?", name).GetErrors(); len(errs) != 0 {
		log.Fatalln(errs)
	}
	return &stu
}

func CheckQQ(qq int64) bool {
	if getStuByQQ(qq).StudentId == 0 {
		return false
	}
	return true
}
