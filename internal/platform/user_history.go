package platform

import (
	"database/sql"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	"github.com/knockbox/authentication/internal/queries"
	"github.com/knockbox/authentication/pkg/models"
	"github.com/knockbox/authentication/pkg/utils"
)

type UserHistorySQLImpl struct {
	*sqlx.DB
	hclog.Logger
}

func (u UserHistorySQLImpl) Create(history models.UserHistory) (sql.Result, error) {
	return utils.Transact(u.DB, func(tx *sql.Tx) (sql.Result, error) {
		return tx.Exec(queries.InsertUserHistory, history.UserId, history.IpAddress, history.Action)
	})
}

func (u UserHistorySQLImpl) GetByUserId(id int, page models.Page) ([]models.UserHistory, error) {
	var history []models.UserHistory
	err := u.Select(history, queries.GetUserHistoryByUserId, id, page.Limit, page.Offset)
	return history, err
}
