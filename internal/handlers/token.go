package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/knockbox/authentication/pkg/keyring"
	"net/http"
)

type Token struct {
	hclog.Logger
	*keyring.KeySet
}

// GetJWKs returns the known public keys that are currently being used for signing.
func (t *Token) GetJWKs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t.GetPublicKeySet())
}

func (t *Token) Route(r *mux.Router) {
	r.HandleFunc("/jwks", t.GetJWKs).Methods(http.MethodGet)
}

func NewToken(l hclog.Logger, ks *keyring.KeySet) *Token {
	return &Token{
		Logger: l,
		KeySet: ks,
	}
}
