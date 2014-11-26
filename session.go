package kitchen

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
)

// SessionContextKey is the key for session store in Context
const SessionContextKey key = 0

type session struct {
	next http.Handler
	sess sessions.Store
}

func (s *session) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, ok := w.(ResponseWriter)
	if !ok {
		logrus.Panic("only usable inside kitchen")
	}

	ctx.SetWithValue(SessionContextKey, s.sess)
	s.next.ServeHTTP(w, r)
}

// SessionMiddlewareGenerator to create new session midleware.
func SessionMiddlewareGenerator(sess sessions.Store) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return &session{next, sess}
	}
}

// GetSessionStore is for get session from context. make sure the SessionMiddleware is there and you call this
// after session middleware in chain
func GetSessionStore(w http.ResponseWriter) (sessions.Store, error) {
	if ctx, ok := w.(ResponseWriter); ok {
		if s, ok := ctx.Context().Value(SessionContextKey).(sessions.Store); ok {
			return s, nil
		}
	}
	return nil, errors.New("you are not inside kitchen or session middleware is not used")
}
