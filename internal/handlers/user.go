package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/knockbox/authentication/internal/client"
	"github.com/knockbox/authentication/pkg/keyring"
	"github.com/knockbox/authentication/pkg/payloads"
	"github.com/knockbox/authentication/pkg/responses"
	"github.com/knockbox/authentication/pkg/utils"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
)

type User struct {
	hclog.Logger
	*keyring.KeySet
	c *client.UserClient
}

func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.UserRegister{}
	if utils.DecodeAndValidateStruct(w, r, payload) {
		return
	}

	if err := u.c.RegisterUser(payload); err != nil {
		if utils.IsDuplicateEntry(err) {
			msg := "a user with the provided username or email already exists"
			responses.NewGenericError(msg).Encode(w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Error(w, "failed to register user", http.StatusInternalServerError)
		u.Error("user registration failed", "error", err, "payload", payload)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.UserLogin{}
	if utils.DecodeAndValidateStruct(w, r, payload) {
		return
	}

	user, err := u.c.GetUserByUsername(payload.Username)
	if err != nil {
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
		responses.NewGenericError("failed to create token").Encode(w)
		w.WriteHeader(http.StatusInternalServerError)

		u.Error("failed to create token", "err", err)
		return
	}

	key := u.GetRandomKey()
	bs, err := jwt.Sign(token, jwt.WithKey(key.Algorithm(), key))
	if err != nil {
		responses.NewGenericError("failed to sign token").Encode(w)
		w.WriteHeader(http.StatusInternalServerError)

		u.Error("failed to sign token", "err", err)
		return
	}

	responses.NewBearerToken(bs, int(u.GetTokenDuration().Seconds())).Encode(w)
}

func (u *User) Route(r *mux.Router) {
	r.HandleFunc("/register", u.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", u.Login).Methods(http.MethodPost)
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
