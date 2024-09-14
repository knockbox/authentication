package client

import (
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	"github.com/knockbox/authentication/internal/platform"
	"github.com/knockbox/authentication/pkg/accessors"
	"github.com/knockbox/authentication/pkg/models"
	"github.com/knockbox/authentication/pkg/payloads"
)

// UserClient provides database functionality for models.User, models.UserDetails and models.UserHistory
type UserClient struct {
	user    accessors.UserAccessor
	details accessors.UserDetailsAccessor
	history accessors.UserHistoryAccessor
	hclog.Logger
}

// NewUserClient creates a new UserClient using the SQLImpl accessors.
func NewUserClient(db *sqlx.DB, l hclog.Logger) *UserClient {
	return &UserClient{
		user: platform.UserSQLImpl{
			DB:     db,
			Logger: l,
		},
		details: platform.UserDetailsSQLImpl{
			DB:     db,
			Logger: l,
		},
		history: platform.UserHistorySQLImpl{
			DB:     db,
			Logger: l,
		},
	}
}

func (c *UserClient) RegisterUser(payload *payloads.UserRegister) error {
	user := models.NewUser()
	if err := user.ApplyRegister(payload); err != nil {
		return err
	}

	result, err := c.user.Create(*user)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if _, err := c.details.CreateForUser(int(id)); err != nil {
		return err
	}

	return err
}

func (c *UserClient) UpdateUser(user *models.User, payload *payloads.UserUpdate) error {
	if err := user.ApplyUpdate(payload); err != nil {
		return err
	}

	_, err := c.user.Update(*user)
	return err
}

func (c *UserClient) GetUserById(id int) (*models.User, error) {
	return c.user.GetById(id)
}

func (c *UserClient) GetUserByAccountId(accountId int) (*models.User, error) {
	return c.user.GetByAccountId(accountId)
}

func (c *UserClient) GetUserByUsername(username string) (*models.User, error) {
	return c.user.GetByUsername(username)
}

func (c *UserClient) GetUsersLikeUsername(username string, page *models.Page) ([]models.User, error) {
	if page == nil || page.Limit == 0 {
		page = models.DefaultPage()
	}

	return c.user.GetLikeUsername(username, *page)
}

func (c *UserClient) DeleteById(id int) error {
	_, err := c.user.DeleteById(id)
	return err
}

func (c *UserClient) UpdateUserDetails(userDetails *models.UserDetails, payload *payloads.UserDetailsUpdate) error {
	userDetails.ApplyUpdate(payload)
	_, err := c.details.Update(*userDetails)
	return err
}
