// This is base on alice. I think alice is good but
// whitout call(next(next(next))), I think this is not a good thing
// to have this kind of chain. also the ResponseWriter need more extra
// Data to handle
package kitchen

import "net/http"

// A simple type, middlewares are like this
type MiddlewareFunc func(http.Handler) http.Handler

// Chain structure for handling middleware
type MiddlewareChain struct {
	functions []MiddlewareFunc
}

// A simple hack middleware to change the ResponseWriter type
func responseWriterWrap(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(NewResponseWriter(rw), r)
	}

	return http.HandlerFunc(fn)
}

// Create the real http handler function.
func (mc MiddlewareChain) Then(h http.Handler) http.Handler {
	var final http.Handler
	if h != nil {
		final = h
	} else {
		final = http.DefaultServeMux
	}
	for i := len(mc.functions) - 1; i >= 0; i-- {
		final = mc.functions[i](final)
	}

	return responseWriterWrap(final)
}

// Create the real http handler function.
func (mc MiddlewareChain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return mc.Then(nil)
	}

	return mc.Then(http.HandlerFunc(fn))
}
