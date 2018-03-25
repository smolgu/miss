package models

import "github.com/go-xorm/xorm"

var (
	dsn = "postgres://postgres:2bcedbe9fcc5fe19568b49a22803a6c9@dokku-postgres-missdb:5432/missdb"
	db  *xorm.Engine
)
