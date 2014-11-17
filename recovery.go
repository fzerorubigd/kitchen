package kitchen

import (
	"net/http"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
)

// RecoveryMiddleware is a simple middleware for prevent program to panic in unhandled
// panic from other middleware or request
// Must add this before any other middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				stack := debug.Stack()
				logrus.WithField("error", err).Warn(err, string(stack))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
