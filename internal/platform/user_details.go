package platform

import (
	"database/sql"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	"github.com/knockbox/authentication/internal/queries"
	"github.com/knockbox/authentication/pkg/models"
	"github.com/knockbox/authentication/pkg/utils"
)

type UserDetailsSQLImpl struct {
	*sqlx.DB
	hclog.Logger
}

func (u UserDetailsSQLImpl) CreateForUser(userId int) (sql.Result, error) {
	return utils.Transact(u.DB, func(tx *sql.Tx) (sql.Result, error) {
		return tx.Exec(queries.InsertUserDetails, userId)
	})
}

func (u UserDetailsSQLImpl) Update(details models.UserDetails) (sql.Result, error) {
	return utils.Transact(u.DB, func(tx *sql.Tx) (sql.Result, error) {
		return tx.Exec(queries.UpdateUserDetails, details.ProfilePicture, details.FullName, details.GithubURL, details.TwitterURL, details.WebsiteURL, details.UserId)
	})
}
