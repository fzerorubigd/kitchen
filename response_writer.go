package kitchen

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

type key string

// ResponseWriter is much like negroni ResponseWriter interface.
// But it also support context, using google context package, you can change the context,
// normally with a child context with old context as its parent
type ResponseWriter interface {
	// The original response writer
	http.ResponseWriter
	// Flusher interface
	http.Flusher
	//The status code set by the client
	Status() int
	// How many byte written into this writer
	Size() int
	// Is the output is written already or not?
	Written() bool
	// Replace context
	SetContext(context.Context)
	// Get the current context
	Context() context.Context
}

// NewResponseWriter Create new response writer base on http response writer interface
// TODO : do i need to implement hijacker and flusher interface?
func NewResponseWriter(rw http.ResponseWriter, ctx context.Context) ResponseWriter {
	return &responseWriter{rw, &sync.RWMutex{}, ctx, 0, 0}
}

type responseWriter struct {
	http.ResponseWriter
	*sync.RWMutex
	ctx    context.Context
	status int
	size   int
}

func (rw *responseWriter) Write(d []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}

	c, err := rw.ResponseWriter.Write(d)
	rw.size += c

	return c, err
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	// I don't know, if this is correct and the http.ResponseWriter in first argument of
	// the http handler always has a CloseNotifier interface or not
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (rw *responseWriter) Written() bool {
	return rw.size > 0
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) SetContext(ctx context.Context) {
	rw.Lock()
	defer rw.Unlock()

	rw.ctx = ctx
}

func (rw *responseWriter) Context() context.Context {
	rw.RLock()
	defer rw.RUnlock()

	return rw.ctx // Is this correct to use this kind of lock here?
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
