package middleware

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/knockbox/authentication/pkg/enums"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"os"
)

var BearerTokenContextKey = "bearer-token"

type BearerToken struct {
	l   hclog.Logger
	set jwk.Set
}

// Middleware is the default handler that rejects with 401 if the token is missing or unverified.
func (b *BearerToken) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseHeader(r.Header, "Authorization", jwt.WithKeySet(b.set))
		if err != nil {
			b.l.Info("bearer token middleware (required)", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := token.PrivateClaims()
		rawRole, ok := claims["role"].(string)
		if !ok {
			b.l.Warn("failed to extract role from claims")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if enums.UserRoleFromString(rawRole).IsForbidden() {
			b.l.Debug("missing required role to access endpoint")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), BearerTokenContextKey, &token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalMiddleware is the non-default handler that allows non-verified or missing tokens. It however
// will only put the token if it is present.
func (b *BearerToken) OptionalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseHeader(r.Header, "Authorization", jwt.WithKeySet(b.set))
		if err != nil {
			ctx := context.WithValue(r.Context(), BearerTokenContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), BearerTokenContextKey, &token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseBearerToken(l hclog.Logger) *BearerToken {
	c := jwk.NewCache(context.Background())
	if err := c.Register(os.Getenv("JWKS_URL")); err != nil {
		panic(err)
	}

	return &BearerToken{
		l:   l,
		set: jwk.NewCachedSet(c, os.Getenv("JWKS_URL")),
	}
}
