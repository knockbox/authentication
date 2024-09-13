package keyring

import "github.com/lestrrat-go/jwx/v2/jwk"

type KeySetResponse struct {
	Keys []jwk.Key `json:"keys"`
}
