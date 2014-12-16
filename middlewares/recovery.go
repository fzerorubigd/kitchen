package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
)

type recoery struct {
	next http.Handler
}

func (re *recoery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			stack := debug.Stack()
			logrus.WithField("error", err).Warn(err, string(stack))
		}
	}()

	re.next.ServeHTTP(w, r)
}

// RecoveryMiddleware is a simple middleware for prevent program to panic in unhandled
// panic from other middleware or request
// Must add this before any other middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	return &recoery{next}
}
