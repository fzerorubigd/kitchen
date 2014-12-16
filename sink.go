package kitchen

import (
	"net/http"

	"golang.org/x/net/context"
)

// Middleware is the call back used for middlewares
type Middleware func(http.Handler) http.Handler

// Chain structure for handling middleware
// This is base on alice. I think alice is good but
// whitout call(next(next(next))), I think this is not a good thing
// to have this kind of chain. also the ResponseWriter need more extra
// Data to handle
type Chain struct {
	functions []Middleware
}

// NewChain create new middleware chain base on provided middlewares function
func NewChain(f ...Middleware) Chain {
	return Chain{f}
}

// A simple hack middleware to change the ResponseWriter type
// The context trigger Context cancel function after request is finished.
func responseWriterWrap(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		// Add support for cancel. when this middleware is done, the request is dead.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		next.ServeHTTP(NewResponseWriter(rw, ctx), r)
	}

	return http.HandlerFunc(fn)
}

// Then Create the real http handler function.
func (mc Chain) Then(h http.Handler) http.Handler {
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
func (mc Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return mc.Then(nil)
	}

	return mc.Then(http.HandlerFunc(fn))
}

// Extend Append function to middleware chain and return NEW chain object
// The old chain is usable after this.
func (mc Chain) Extend(f ...Middleware) Chain {
	newFuncs := make([]Middleware, len(mc.functions))
	copy(newFuncs, mc.functions)
	newFuncs = append(newFuncs, f...)

	newChain := NewChain(newFuncs...)
	return newChain
}
