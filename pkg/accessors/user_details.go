package accessors

import (
	"database/sql"
	"github.com/knockbox/authentication/pkg/models"
)

// UserDetailsAccessor defines all queries available for models.UserDetails
type UserDetailsAccessor interface {
	CreateForUser(userId int) (sql.Result, error)
	Update(details models.UserDetails) (sql.Result, error)
}
