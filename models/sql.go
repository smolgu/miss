package models

import (
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

//var db *xorm.Enginek

func NewContext() (err error) {
	db, err = xorm.NewEngine("postgres", "postgres://postgres:postgres@127.0.0.1:5432/missdb?sslmode=disable")
	if err != nil {
		return
	}

	db.ShowSQL()

	err = db.Sync2(&userVoteModel{})

	return
}
