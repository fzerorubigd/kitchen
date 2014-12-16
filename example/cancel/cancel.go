package main

import (
	"log"
	"net/http"
	"time"

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

	start := time.Now()
	log.Println("New request...")

	go func() {
		<-krw.Context().Done()
		latency := time.Since(start)
		log.Println(r.URL.Path, krw.Status(), krw.Context().Err(), latency)
	}()

	m.next.ServeHTTP(w, r) // Call the next one
}

func newExampleMiddleware(next http.Handler) http.Handler {
	return &midleware{next}
}

func main() {

	req := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("See terminal for log."))
	}

	http.Handle(
		"/",
		kitchen.NewChain(newExampleMiddleware).ThenFunc(req),
	)

	log.Fatal(http.ListenAndServe(":9091", nil))
}
