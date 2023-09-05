package middleware

import (
	"context"
	"net/http"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var contextKeyUser = contextKey("user")

func BasicAuthMiddleware(next http.Handler) http.Handler {
	users := map[string]string{
		"aandmelom@falabella.cl":           "123456",
		"jhernandezme@abcservicios.com.co": "123456",
		"malexlondono@abcservicios.com.co": "123456",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		pass, exist := users[user]
		if !exist {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if password != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyUser, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
