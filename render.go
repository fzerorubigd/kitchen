package kitchen

import (
	"errors"
	"net/http"

	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

// RenderContextKey is the key used to store render inside the context
const RenderContextKey key = 1

// RenderMiddlewareGenerator generate a new render middleware for use in kitchen using render package
// Personally I hate when the framework automatically render a template base on its name. so its
// not an option here.
func RenderMiddlewareGenerator(r *render.Render) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Do not call panic, here. I think its ok to just ignore this requests
				if ctx, ok := w.(ResponseWriter); ok {
					ctx.SetContext(context.WithValue(ctx.Context(), RenderContextKey, r))
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}

// GetRender is a helper function and returns the render object.
// Reterns error in case of wrong interface or when the middleware is not used on the request.
func GetRender(w http.ResponseWriter) (*render.Render, error) {
	ctx, ok := w.(ResponseWriter)
	if !ok {
		return nil, errors.New("the Context interface is not implemented")
	}
	r, ok := ctx.Context().Value(RenderContextKey).(*render.Render)
	if ok {
		return r, nil
	}

	return nil, errors.New("the render is not stored here")
}
