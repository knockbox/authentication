package platform

import (
	"database/sql"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	"github.com/knockbox/authentication/internal/queries"
	"github.com/knockbox/authentication/pkg/models"
	"github.com/knockbox/authentication/pkg/utils"
)

type UserSQLImpl struct {
	*sqlx.DB
	hclog.Logger
}

func (u UserSQLImpl) Create(user models.User) (sql.Result, error) {
	return utils.Transact(u.DB, func(tx *sql.Tx) (sql.Result, error) {
		return tx.Exec(queries.InsertUser, user.AccountId, user.Username, user.Password, user.Email, user.Role)
	})
}

func (u UserSQLImpl) Update(user models.User) (sql.Result, error) {
	return utils.Transact(u.DB, func(tx *sql.Tx) (sql.Result, error) {
		return tx.Exec(queries.UpdateUser, user.Email, user.Password, user.Role, user.Id)
	})
}

func (u UserSQLImpl) GetById(id int) (*models.User, error) {
	user := &models.User{}
	err := u.Get(user, queries.GetUserById, id)
	return user, err
}

func (u UserSQLImpl) GetByAccountId(id int) (*models.User, error) {
	user := &models.User{}
	err := u.Get(user, queries.GetUserByAccountId, id)
	return user, err
}

func (u UserSQLImpl) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := u.Get(user, queries.GetUserByUsername, username)
	return user, err
}

func (u UserSQLImpl) GetLikeUsername(username string, page models.Page) ([]models.User, error) {
	var users []models.User
	err := u.Select(users, queries.GetUsersLikeUsername, username, page.Limit, page.Offset)
	return users, err
}

func (u UserSQLImpl) DeleteById(id int) (sql.Result, error) {
	return utils.Transact(u.DB, func(tx *sql.Tx) (sql.Result, error) {
		return tx.Exec(queries.DeleteUserById, id)
	})
}
