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
	leaves, err := QueryLeaveByID(21011631)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(leaves)
}

func TestQueryTodayLeaves(t *testing.T) {
	initTest(t)
	leaves, err := QueryTodayLeaves()
	if err != nil {
		t.Error(err)
	}
	t.Log(leaves)
}

func TestQueryShouldArrive(t *testing.T) {
	initTest(t)
	list, err := QueryShouldArrive()
	t.Log(list)
	if err != nil {
		t.Error(err)
	}
}
