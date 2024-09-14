package accessors

import (
	"database/sql"
	"github.com/knockbox/authentication/pkg/models"
)

// UserHistoryAccessor defines all queries available for models.UserHistory
type UserHistoryAccessor interface {
	Create(history models.UserHistory) (sql.Result, error)
	GetByUserId(id int, page models.Page) ([]models.UserHistory, error)
}
