package kitchen

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

func TestResponseWriterWritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, context.Background())

	rw.Write([]byte("Hello world"))

	assert.Equal(t, rec.Code, rw.Status())
	assert.Equal(t, rec.Body.String(), "Hello world")
	assert.Equal(t, rw.Status(), http.StatusOK)
	assert.Equal(t, rw.Size(), 11)
	assert.Equal(t, rw.Written(), true)
}

func TestResponseWriterWritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, context.Background())

	rw.Write([]byte("Hello world"))
	rw.Write([]byte("foo bar bat baz"))

	assert.Equal(t, rec.Code, rw.Status())
	assert.Equal(t, rec.Body.String(), "Hello worldfoo bar bat baz")
	assert.Equal(t, rw.Status(), http.StatusOK)
	assert.Equal(t, rw.Size(), 26)
}

func TestResponseWriterWritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, context.Background())

	rw.WriteHeader(http.StatusNotFound)

	assert.Equal(t, rec.Code, rw.Status())
	assert.Equal(t, rec.Body.String(), "")
	assert.Equal(t, rw.Status(), http.StatusNotFound)
	assert.Equal(t, rw.Size(), 0)
}

func TestResponseWriterCloseNotify(t *testing.T) {
	rec := newCloseNotifyingRecorder()
	rw := NewResponseWriter(rec, context.Background())
	closed := false
	notifier := rw.(http.CloseNotifier).CloseNotify()
	rec.close()
	select {
	case <-notifier:
		closed = true
	case <-time.After(time.Second):
	}
	assert.Equal(t, closed, true)
}

func TestResponseWriterFlusher(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec, context.Background())

	f, ok := rw.(http.Flusher)
	f.Flush()
	assert.Equal(t, ok, true)
}

func TestContextRelated(t *testing.T) {
	rec := httptest.NewRecorder()
	ctx := context.Background()
	rw := NewResponseWriter(rec, ctx)

	assert.Equal(t, ctx, rw.Context())

	rw.SetWithDeadline(time.Now())
	rw.SetWithTimeout(time.Hour)
	rw.SetWithValue("example", "data")

	d, ok := rw.Context().Value("example").(string)

	assert.Equal(t, true, ok)
	assert.Equal(t, "data", d)

}
