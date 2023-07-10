package utils

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrUniqueKeyViolation = errors.New("unique key violation")
	ErrRecordNotFound     = errors.New("record not found")
)

func TranslateErrors(dialect string, e error) error {
	switch dialect {
	case "sqlite":
		if err, ok := e.(sqlite3.Error); ok {
			if err.ExtendedCode == sqlite3.ErrConstraintUnique {
				return ErrUniqueKeyViolation
			}

			if err.Code == sqlite3.ErrNotFound {
				return ErrRecordNotFound
			}
			// other errors handling from sqlite3
		}
	case "postgres":
		if err, ok := e.(*pgconn.PgError); ok {
			if err.Code == pgerrcode.UniqueViolation {
				return ErrUniqueKeyViolation
			}

			if err.Code == pgerrcode.CaseNotFound {
				return ErrRecordNotFound
			}

			// other errors handling from pgconn
		}
	}
	return nil
}
