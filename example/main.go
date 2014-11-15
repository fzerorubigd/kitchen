package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fzerorubigd/kitchen"
)

func text(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
}

func main() {
	http.Handle("/", kitchen.NewMiddlewareChain(kitchen.RecoveryMiddleware, kitchen.LoggerMiddleware).ThenFunc(text))

	log.Fatal(http.ListenAndServe(":9091", nil))
}
