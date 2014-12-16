package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/kitchen"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
)

type Controller struct{}

func (c Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, ok := w.(context.Context)
	if !ok {
		logrus.Panic("This is not inside kitchen")
	}

	ctx.Value("test")
}

func test(w http.ResponseWriter, r *http.Request) {
	go func() {
		d := w.(kitchen.ResponseWriter).Context().Done()
		fmt.Print("Waiting...")
		<-d
		fmt.Print("Im done!")
	}()

	if store, err := kitchen.GetSessionStore(w); err == nil {
		if sess, err := store.Get(r, "test-store"); err == nil {

			if x, ok := sess.Values["data"]; ok {
				fmt.Fprintf(w, x.(string))
			} else {
				// DO NOT SEND DATA TO OUTPUT BEFORE SAVE SESSION!
				sess.Values["data"] = "abcd"
				sess.Save(r, w)
				fmt.Fprintf(w, "Setting session for the first time")
			}

		} else {
			logrus.Panic(err)
		}
	} else {
		logrus.Panic(err)
	}
	//d := w.(kitchen.ResponseWriter).Context().Done()
	//<-d
	fmt.Fprintf(w, ":)")
}

func main() {
	store := sessions.NewCookieStore([]byte("something...."))
	http.Handle(
		"/",
		kitchen.NewChain(
			kitchen.RecoveryMiddleware,
			kitchen.TimeoutMiddlewareGenerator(time.Second*10),
			kitchen.LoggerMiddleware,
			kitchen.SessionMiddlewareGenerator(store),
		).ThenFunc(test),
	)

	log.Fatal(http.ListenAndServe(":9091", nil))
}
