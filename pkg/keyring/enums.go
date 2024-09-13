package keyring

import (
	"crypto/elliptic"
	"github.com/lestrrat-go/jwx/v2/jwa"
)

type CurveType string

const (
	P256 CurveType = "P-256"
	P384           = "P-384"
	P521           = "P-521"
)

func (c CurveType) GetEllipticCurve() elliptic.Curve {
	switch c {
	case P256:
		return elliptic.P256()
	case P384:
		return elliptic.P384()
	case P521:
		return elliptic.P521()
	}

	return elliptic.P521()
}

func (c CurveType) GetAlgorithm() jwa.SignatureAlgorithm {
	switch c {
	case P256:
		return jwa.ES256
	case P384:
		return jwa.ES384
	case P521:
		return jwa.ES512
	}

	return jwa.ES512
}
