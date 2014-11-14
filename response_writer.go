package kitchen

import "net/http"

// negroni use the http.Flusher interface too. do I need this?
type ResponseWriter interface {
	http.ResponseWriter
	//The status code set by the client
	Status() int
	// How many byte written into this writer
	Size() int
	// Is the output is written already or not?
	Written() bool
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0, 0}
}

type responseWriter struct {
	http.ResponseWriter
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
