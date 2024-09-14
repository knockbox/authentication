package queries

import _ "embed"

//go:embed user/insert.sql
var InsertUser string

//go:embed user/update.sql
var UpdateUser string

//go:embed user/select-by-id.sql
var GetUserById string

//go:embed user/select-by-account_id.sql
var GetUserByAccountId string

//go:embed user/select-by-username.sql
var GetUserByUsername string

//go:embed user/like-username.sql
var GetUsersLikeUsername string

//go:embed user/delete-by-id.sql
var DeleteUserById string
