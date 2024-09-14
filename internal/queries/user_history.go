package queries

import _ "embed"

//go:embed user-history/insert.sql
var InsertUserHistory string

//go:embed user-history/select-by-user_id.sql
var GetUserHistoryByUserId string
