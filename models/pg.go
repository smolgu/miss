package models

import "github.com/go-xorm/xorm"

var (
	dsn = "postgres://postgres:bbfedc6d1142829c0c07e68052f79b3c@dokku-postgres-smolgudb:5432/smolgudb"
	db  *xorm.Engine
)
