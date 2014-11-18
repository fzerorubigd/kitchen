package kitchen

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// TimeoutMiddlewareGenerator generate a middleware for trigger Done channel on context
// in every request
func TimeoutMiddlewareGenerator(timeout time.Duration) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if ctx, ok := w.(ResponseWriter); ok {
					// We call cancel on parent and so this cancel is not required (I THINK :) )
					c, _ := context.WithTimeout(ctx.Context(), timeout)
					ctx.SetContext(c)
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
