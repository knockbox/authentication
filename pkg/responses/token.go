package responses

import (
	"encoding/json"
	"net/http"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewBearerToken(accessToken []byte, expiresInSeconds int) *Token {
	return &Token{
		AccessToken: string(accessToken),
		TokenType:   "Bearer",
		ExpiresIn:   expiresInSeconds,
	}
}

func (t *Token) Encode(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}
