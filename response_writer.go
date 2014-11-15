package kitchen

import (
	"net/http"

	"golang.org/x/net/context"
)

// Much like negroni ResponseWriter interface.
// But it also support context, using google context package, you can change the context,
// normally with a child context with old context as its parent
type ResponseWriter interface {
	// The original response writer
	http.ResponseWriter
	// Flusher interface
	http.Flusher
	// This is also Context support
	context.Context
	//The status code set by the client
	Status() int
	// How many byte written into this writer
	Size() int
	// Is the output is written already or not?
	Written() bool
	// Replace context, need to think more :)
	SetContext(context.Context)
}

// Create new response writer base on http response writer interface
// TODO : do i need to implement hijacker and flusher interface?
func NewResponseWriter(rw http.ResponseWriter, ctx context.Context) ResponseWriter {
	return &responseWriter{rw, ctx, 0, 0}
}

type responseWriter struct {
	http.ResponseWriter
	context.Context
	status int
	size   int
}

func (rw responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
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
	rw.Context = ctx
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
