package data

import (
	"github.com/BurntSushi/toml"
	"testing"
)

func initTest(t *testing.T) {
	conf := make(map[string]string)
	if _, err := toml.DecodeFile("./../config.toml", &conf); err != nil {
		t.Error(err)
	}
	if err := InitDB(conf["dbsource"]); err != nil {
		t.Error(err)
	}
}

func TestQueryLeaveByID(t *testing.T) {
	initTest(t)
	leaves, err := QueryLeaveByID(21181305)
	if err != nil {
		t.Error(err)
	}
	t.Log(leaves[0])
}

func TestQueryTodayLeaves(t *testing.T) {
	initTest(t)
	leaves, err := QueryTodayLeaves()
	if err != nil {
		t.Error(err)
	}
	t.Log(leaves)
}

func TestIsTodayLeaveByID(t *testing.T) {
	initTest(t)
	yes, err := IsTodayLeaveByID(21181305)
	if err != nil {
		t.Error(err)
	}
	t.Log(yes)
}
