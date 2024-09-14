package utils

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// Transact wraps sql transaction rollback and commit functionality.
func Transact(db *sqlx.DB, txFunc func(tx *sql.Tx) (sql.Result, error)) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			// If we panic, re-throw but rollback first.
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// Rollback if we errored.
			_ = tx.Rollback()
		} else {
			// Commit it.
			err = tx.Commit()
		}
	}()

	result, err := txFunc(tx)
	return result, err
}
