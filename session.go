package kitchen

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
)

// SessionContextKey is the key for session store in Context
const SessionContextKey key = 0

// SessionMiddleware create new session midleware.
// TODO : is this idiomatic :) to have 3 nested function like this?
func SessionMiddlewareGenerator(sess sessions.Store) MiddlewareFunc {
	fn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, ok := w.(ResponseWriter)
			if !ok {
				logrus.Panic("only usable inside kitchen")
			}

			ctx.SetContext(context.WithValue(ctx.Context(), SessionContextKey, sess))
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return fn
}

// GetSessionStore is for get session from context. make sure the SessionMiddleware is there and you call this
// after session middleware in chain
func GetSessionStore(w http.ResponseWriter) (sessions.Store, error) {
	ctx, ok := w.(ResponseWriter)
	if !ok {
		return nil, errors.New("the Context interface is not implemented")
	}
	s, ok := ctx.Context().Value(SessionContextKey).(sessions.Store)
	if ok {
		return s, nil
	}

	return nil, errors.New("the session is not stored here")
}
