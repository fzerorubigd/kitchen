package kitchen

import (
	"net/http"

	"golang.org/x/net/context"
)

// MiddlewareFunc type, middlewares are like this
type MiddlewareFunc func(http.Handler) http.Handler

// MiddlewareChain structure for handling middleware
// This is base on alice. I think alice is good but
// whitout call(next(next(next))), I think this is not a good thing
// to have this kind of chain. also the ResponseWriter need more extra
// Data to handle
type MiddlewareChain struct {
	functions []MiddlewareFunc
}

// NewMiddlewareChain create new middleware base on functions
func NewMiddlewareChain(f ...MiddlewareFunc) MiddlewareChain {
	return MiddlewareChain{f}
}

// A simple hack middleware to change the ResponseWriter type
func responseWriterWrap(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.TODO() // :)
		next.ServeHTTP(NewResponseWriter(rw, ctx), r)
	}

	return http.HandlerFunc(fn)
}

// Then Create the real http handler function.
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

// ThenFunc Create the real http handler function.
func (mc MiddlewareChain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return mc.Then(nil)
	}

	return mc.Then(http.HandlerFunc(fn))
}

// Extend Append function to middleware chain and return NEW chain object
// The old chain is usable after this.
func (mc MiddlewareChain) Extend(f ...MiddlewareFunc) MiddlewareChain {
	newFuncs := make([]MiddlewareFunc, len(mc.functions))
	copy(newFuncs, mc.functions)
	newFuncs = append(newFuncs, f...)

	newChain := NewMiddlewareChain(newFuncs...)
	return newChain
}
