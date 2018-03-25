package models

import (
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	// pq postgresql database driver
	_ "github.com/lib/pq"

	"github.com/smolgu/miss/pkg/setting"
)

// NewContext open connection to database, migrate schema
func NewContext() (err error) {
	if setting.Dev {
		db, err = xorm.NewEngine("postgres", "postgres://postgres:postgres@127.0.0.1:5432/missdb?sslmode=disable")
		if err != nil {
			return errors.Wrap(err, "open database connection")
		}
		db.ShowSQL()
	} else {
		db, err = xorm.NewEngine("postgres", dsn)
		if err != nil {
			return errors.Wrap(err, "open database connection")
		}
	}

	err = db.Sync2(&userVoteModel{})
	if err != nil {
		return errors.Wrap(err, "sync2 (migration)")
	}

	return
}
