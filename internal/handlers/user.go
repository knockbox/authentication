package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/knockbox/authentication/internal/client"
	"github.com/knockbox/authentication/pkg/payloads"
	"github.com/knockbox/authentication/pkg/responses"
	"github.com/knockbox/authentication/pkg/utils"
	"net/http"
)

type User struct {
	hclog.Logger
	c *client.UserClient
}

func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.UserRegister{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, "malformed request body, expected json", http.StatusBadRequest)
		return
	}

	if errs := utils.ValidateStruct(payload); errs != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errs)
		return
	}

	if err := u.c.RegisterUser(payload); err != nil {
		if utils.IsDuplicateEntry(err) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			ge := responses.NewGenericError("a user with the provided username or email already exists")
			_ = json.NewEncoder(w).Encode(ge)
			return
		}

		http.Error(w, "failed to register user", http.StatusInternalServerError)
		u.Error("user registration failed", "error", err, "payload", payload)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *User) Route(r *mux.Router) {
	r.HandleFunc("/register", u.Register).Methods(http.MethodPost)
}

func NewUser(l hclog.Logger) *User {
	db, err := utils.MySQLConnection()
	if err != nil {
		panic(err)
	}

	return &User{
		Logger: l,
		c:      client.NewUserClient(db, l),
	}
}
