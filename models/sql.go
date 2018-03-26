package models

import (
	"log"

	"github.com/go-xorm/xorm"
	pkgErrors "github.com/pkg/errors"
	// pq postgresql database driver
	_ "github.com/lib/pq"

	"github.com/smolgu/miss/pkg/errors"
	"github.com/smolgu/miss/pkg/setting"
)

// NewContext open connection to database, migrate schema
func NewContext() (err error) {
	if setting.Dev {
		db, err = xorm.NewEngine("postgres", "postgres://postgres:postgres@127.0.0.1:5432/missdb?sslmode=disable")
		if err != nil {
			return pkgErrors.Wrap(err, "open database connection")
		}
		db.ShowSQL()
	} else {
		db, err = xorm.NewEngine("postgres", dsn)
		if err != nil {
			return pkgErrors.Wrap(err, "open database connection")
		}
	}

	err = db.Sync2(
		new(userVoteModel),
		new(User),
	)
	if err != nil {
		return pkgErrors.Wrap(err, "sync2 (migration)")
	}

	err = checkInstall()
	if err != nil {
		return pkgErrors.Wrap(err, "create admin user")
	}

	return
}

func checkInstall() error {
	_, err := Users.Get(1)
	if err != nil {
		if !errors.CheckTyped(err, errors.ErrNotFound) {
			return err
		}
	}
	user := &User{
		FirstName:       "Администратор",
		MessagesFromAll: true,
	}
	_, err = db.InsertOne(user)
	if err != nil {
		return pkgErrors.Wrap(err, "cannot create admin user")
	}

	token, err := Sessions.New(1)
	if err != nil {
		return pkgErrors.Wrap(err, "cannot create admin jwt token")
	}
	log.Printf("admin token: %v", token)

	userID, err := Sessions.Check(token)
	if err != nil {
		return pkgErrors.Wrap(err, "session check admin token")
	}

	if userID != 1 {
		return pkgErrors.Wrapf(err, "user id in session not equal 1, got: %d", userID)
	}

	return nil
}
