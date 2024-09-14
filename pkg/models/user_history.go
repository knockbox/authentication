package models

import (
	"github.com/knockbox/authentication/pkg/enums"
	"time"
)

// UserHistory defines a User(s) history in our database.db
type UserHistory struct {
	Id        uint             `db:"id"`
	UserId    uint             `db:"user_id"`
	IpAddress string           `db:"ip_address"`
	Timestamp time.Time        `db:"timestamp"`
	Action    enums.UserAction `db:"action"`
}
