package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/knockbox/authentication/internal/client"
	"github.com/knockbox/authentication/pkg/keyring"
	"github.com/knockbox/authentication/pkg/middleware"
	"github.com/knockbox/authentication/pkg/models"
	"github.com/knockbox/authentication/pkg/payloads"
	"github.com/knockbox/authentication/pkg/responses"
	"github.com/knockbox/authentication/pkg/utils"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"strings"
)

type User struct {
	hclog.Logger
	*keyring.KeySet
	c *client.UserClient
}

// Register handles user registration.
func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.UserRegister{}
	if utils.DecodeAndValidateStruct(w, r, payload) {
		return
	}

	if err := u.c.RegisterUser(payload); err != nil {
		if utils.IsDuplicateEntry(err) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			msg := "a user with the provided username or email already exists"
			responses.NewGenericError(msg).Encode(w)
			return
		}

		http.Error(w, "failed to register user", http.StatusInternalServerError)
		u.Error("user registration failed", "error", err, "payload", payload)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles user login.
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.UserLogin{}
	if utils.DecodeAndValidateStruct(w, r, payload) {
		return
	}

	user, err := u.c.GetUserByUsername(payload.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		u.Debug("username not found", "payload", payload, "err", err)
		return
	}

	if !utils.ComparePasswords(user.Password, payload.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := user.CreateToken(u.GetTokenDuration())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		responses.NewGenericError("failed to create token").Encode(w)

		u.Error("failed to create token", "err", err)
		return
	}

	key := u.GetRandomKey()
	bs, err := jwt.Sign(token, jwt.WithKey(key.Algorithm(), key))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		responses.NewGenericError("failed to sign token").Encode(w)

		u.Error("failed to sign token", "err", err)
		return
	}

	responses.NewBearerToken(bs, int(u.GetTokenDuration().Seconds())).Encode(w)
}

// GetByAccountId returns a user by their account_id.
func (u *User) GetByAccountId(w http.ResponseWriter, r *http.Request) {
	accountId, ok := mux.Vars(r)["account_id"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responses.NewGenericError("account_id was not provided").Encode(w)
		return
	}

	if _, err := uuid.Parse(accountId); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responses.NewGenericError("the provided account_id failed to parse").Encode(w)
		return
	}

	user, err := u.c.GetUserByAccountId(accountId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Error("failed to get user by account_id", "err", err)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user.DTO())
}

// GetByUsername returns a user by their username.
func (u *User) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username, ok := mux.Vars(r)["username"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responses.NewGenericError("username was not provided").Encode(w)
		return
	}

	username = strings.TrimSpace(username)
	if len(username) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responses.NewGenericError("provided username was empty").Encode(w)
		return
	}

	user, err := u.c.GetUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Error("failed to get user by username", "err", err)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user.DTO())
}

// GetLikeUsername returns all the users like the given username.
func (u *User) GetLikeUsername(w http.ResponseWriter, r *http.Request) {
	username, ok := mux.Vars(r)["username"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responses.NewGenericError("username was not provided").Encode(w)
		return
	}

	username = strings.TrimSpace(username)
	if len(username) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responses.NewGenericError("provided username was empty").Encode(w)
		return
	}

	page := models.PageFromRequest(r)
	users, err := u.c.GetUsersLikeUsername(username, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Error("failed to get users like username", "err", err)
		return
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// TODO: Return the users and the paging result.
	// TODO: Get the total.
	var dtos []*models.UserDTO
	for _, user := range users {
		dtos = append(dtos, user.DTO())
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dtos)
}

// Update applies changes to the User based on the bearer token.
func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.UserUpdate{}
	if utils.DecodeAndValidateStruct(w, r, payload) {
		return
	}

	if !utils.PayloadHasChanges(*payload) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, ok := r.Context().Value(middleware.BearerTokenContextKey).(*jwt.Token)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		u.Warn("User.Update token was expected and should have existed but was not found")
		return
	}

	accountId, ok := (*token).Get("account_id")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		u.Warn("User.Update token was missing claim 'account_id'")
		return
	}

	id, ok := accountId.(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		u.Warn("User.Update token claim 'account_id' was not a string")
		return
	}

	user, err := u.c.GetUserByAccountId(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Error("failed to get user by account_id", "err", err)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := u.c.UpdateUser(user, payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *User) Route(r *mux.Router) {
	bearer := middleware.UseBearerToken(u.Logger)

	r.HandleFunc("/register", u.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", u.Login).Methods(http.MethodPost)

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/{account_id}", u.GetByAccountId).Methods(http.MethodGet)
	userRouter.HandleFunc("/username/{username}", u.GetByUsername).Methods(http.MethodGet)

	searchRouter := userRouter.PathPrefix("/search").Subrouter()
	searchRouter.HandleFunc("/{username}", u.GetLikeUsername).Methods(http.MethodGet)

	authorizedUserRouter := r.PathPrefix("/user").Subrouter()
	authorizedUserRouter.Use(bearer.Middleware)
	authorizedUserRouter.HandleFunc("", u.Update).Methods(http.MethodPut, http.MethodPatch)
}

func NewUser(l hclog.Logger, ks *keyring.KeySet) *User {
	db, err := utils.MySQLConnection()
	if err != nil {
		panic(err)
	}

	return &User{
		Logger: l,
		KeySet: ks,
		c:      client.NewUserClient(db, l),
	}
}
