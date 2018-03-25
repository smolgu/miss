package models

import (
	"github.com/asdine/storm"
)

var (
	stormDB *storm.DB
)

// NewContext инизиализирует базу данных
// func NewContext() (err error) {
// 	stormDB, err = storm.Open(setting.DataPath)
// 	return
// }
