package accessors

import (
	"database/sql"
	"github.com/knockbox/authentication/pkg/models"
)

// UserAccessor defines all queries available for models.User
type UserAccessor interface {
	Create(user models.User) (sql.Result, error)
	Update(user models.User) (sql.Result, error)
	GetById(id int) (*models.User, error)
	GetByAccountId(accountId string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetLikeUsername(username string, page models.Page) ([]models.User, error)
	DeleteById(id int) (sql.Result, error)
}
