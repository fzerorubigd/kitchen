package kitchen

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type timeout struct {
	next    http.Handler
	timeout time.Duration
}

func (t *timeout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ctx, ok := w.(ResponseWriter); ok {
		// We call cancel on parent and so this cancel is not required (I THINK :) )
		c, _ := context.WithTimeout(ctx.Context(), t.timeout)
		ctx.SetContext(c)
	}

	t.next.ServeHTTP(w, r)

}

// TimeoutMiddlewareGenerator generate a middleware for trigger Done channel on context
// in every request
func TimeoutMiddlewareGenerator(to time.Duration) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return &timeout{next, to}
	}
}
