package middlewares

import (
	"net/http"
	"time"

	"github.com/fzerorubigd/kitchen"
)

type timeout struct {
	next    http.Handler
	timeout time.Duration
}

func (t *timeout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ctx, ok := w.(kitchen.ResponseWriter); ok {
		// We call cancel on parent and so this cancel is not required (I THINK :) )
		ctx.SetWithTimeout(t.timeout)
	}

	t.next.ServeHTTP(w, r)

}

// TimeoutMiddlewareGenerator generate a middleware for trigger Done channel on context
// in every request
func TimeoutMiddlewareGenerator(to time.Duration) kitchen.Middleware {
	return func(next http.Handler) http.Handler {
		return &timeout{next, to}
	}
}
