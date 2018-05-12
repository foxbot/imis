package main

import (
	"net/http"
)

// AuthorizationMiddleware builds a middleware that expects an Authorization header
func AuthorizationMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := r.Header.Get("Authorization")

			if t != token {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
