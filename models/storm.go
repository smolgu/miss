package models

import (
	"github.com/asdine/storm"
	"github.com/smolgu/miss/pkg/setting"
)

var (
	stormDB *storm.DB
)

// NewContext инизиализирует базу данных
func NewContext() (err error) {
	stormDB, err = storm.Open(setting.DataPath)
	return
}
