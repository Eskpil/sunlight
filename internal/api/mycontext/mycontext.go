package mycontext

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Context struct {
	Db *gorm.DB
}

func New() (*Context, error) {
	m := new(Context)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	m.Db = db

	if err := m.Db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}

	return m, nil
}
