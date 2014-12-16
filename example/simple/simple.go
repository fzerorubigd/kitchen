package main

import (
	"log"
	"net/http"

	"github.com/fzerorubigd/kitchen"
)

type midleware struct {
	next http.Handler
}

func (m *midleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	krw, ok := w.(kitchen.ResponseWriter)
	if !ok {
		panic("wtf?")
	}
	krw.SetWithValue("test", "from middleware!")

	// Call the next middleware
	m.next.ServeHTTP(w, r)
}

func newExampleMiddleware(next http.Handler) http.Handler {
	return &midleware{next}
}

func main() {

	req := func(w http.ResponseWriter, r *http.Request) {
		krw, ok := w.(kitchen.ResponseWriter)
		if !ok {
			panic("wtf?")
		}
		test, ok := krw.Context().Value("test").(string)
		if !ok {
			panic("wtf?")
		}

		w.Write([]byte(test))
	}

	http.Handle(
		"/",
		kitchen.NewChain(newExampleMiddleware).ThenFunc(req),
	)

	log.Fatal(http.ListenAndServe(":9091", nil))
}
