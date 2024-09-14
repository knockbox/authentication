package utils

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

// IsDuplicateEntry determines if the given error is a mysql.MySQLError and the error number is 1062.
func IsDuplicateEntry(in error) bool {
	err := &mysql.MySQLError{}
	if ok := errors.As(in, &err); !ok {
		return false
	}

	return err.Number == 1062
}
