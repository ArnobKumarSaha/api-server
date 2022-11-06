package auth

import (
	"context"
	"github.com/Arnobkumarsaha/new-server/schemas"
	"net/http"
)

type AuthProductResource struct {
	*schemas.ProductResource
}

func (rs *AuthProductResource) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tknStr, err := getToken(w, r)
		if err != nil {
			return
		}
		_, err = getClaimFromTokenString(w, tknStr)
		if err != nil {
			return
		}

		// Set key-value pair in the context, before the next call
		ctx := context.WithValue(r.Context(), "token", tknStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
