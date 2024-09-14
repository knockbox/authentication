package queries

import _ "embed"

//go:embed user-details/insert.sql
var InsertUserDetails string

//go:embed user-details/update.sql
var UpdateUserDetails string
