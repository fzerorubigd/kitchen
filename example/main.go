package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/kitchen"
	"github.com/gorilla/sessions"
)

func text(w http.ResponseWriter, r *http.Request) {
	if store, err := kitchen.GetSessionStore(w); err == nil {
		if sess, err := store.Get(r, "test-store"); err == nil {
			//sess.Values["test"] = "abcd"
			//sess.Save(r, w)
			fmt.Fprintf(w, sess.Values["test"].(string))

		} else {
			logrus.Panic(err)
		}
	} else {
		logrus.Panic(err)
	}
	fmt.Fprintf(w, "Hi")
}

func main() {
	store := sessions.NewCookieStore([]byte("something...."))
	http.Handle(
		"/",
		kitchen.NewMiddlewareChain(
			kitchen.RecoveryMiddleware,
			kitchen.LoggerMiddleware,
			kitchen.SessionMiddleware(store),
		).ThenFunc(text),
	)

	log.Fatal(http.ListenAndServe(":9091", nil))
}
