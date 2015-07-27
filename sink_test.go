package kitchen

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {

	var f1, f2, f3 Middleware

	f1 = func(http.Handler) http.Handler {
		return nil
	}

	f2 = func(http.Handler) http.Handler {
		return nil
	}

	f3 = func(http.Handler) http.Handler {
		return nil
	}

	chain := NewChain(f1, f2)

	assert.Equal(t, reflect.ValueOf(chain.functions[0]), reflect.ValueOf(f1))
	assert.Equal(t, reflect.ValueOf(chain.functions[1]), reflect.ValueOf(f2))

	assert.Len(t, chain.functions, 2)
	cpy := chain.Extend(f3)

	assert.Equal(t, reflect.ValueOf(cpy.functions[0]), reflect.ValueOf(f1))
	assert.Equal(t, reflect.ValueOf(cpy.functions[1]), reflect.ValueOf(f2))
	assert.Equal(t, reflect.ValueOf(cpy.functions[2]), reflect.ValueOf(f3))
	assert.Len(t, chain.functions, 2)
	assert.Len(t, cpy.functions, 3)
}

func TestChainThen(t *testing.T) {

	var counter int = 0

	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, ok := w.(ResponseWriter)
				assert.Equal(t, true, ok)
				counter++

				next.ServeHTTP(w, r)
			},
		)
	}

	last := func(w http.ResponseWriter, r *http.Request) {
		_, ok := w.(ResponseWriter)
		assert.Equal(t, true, ok)
		counter++
	}

	chain := NewChain(f, f, f)
	handler := chain.ThenFunc(last)
	handler2 := chain.ThenFunc(nil)
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, 4, counter)

	counter = 0
	handler2.ServeHTTP(w, req)
	assert.Equal(t, 3, counter)
}
