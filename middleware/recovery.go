package middleware

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

// A simple middleware for prevent program to panic in unhandled
// panic from other middleware or request
// Must add this before any other middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Warn(r)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
