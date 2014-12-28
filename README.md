# kitchen

[![Build Status](https://travis-ci.org/fzerorubigd/kitchen.svg?branch=master)](https://travis-ci.org/fzerorubigd/kitchen)
[![Coverage Status](https://coveralls.io/repos/fzerorubigd/kitchen/badge.png)](https://coveralls.io/r/fzerorubigd/kitchen)

A very simple bootstrap for my web project, based on [alice](https://github.com/justinas/alice) and [negroni](https://github.com/codegangsta/negroni) and [context]() package.

## How?

Its very simple to handle a request in `net/http` package, and kitchen just add a simple middleware system to that, like this :

```go
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
```

the middleware system, always add a responseWriterWrap middleware at the begining, so the first argument is always `kitchen.ResponseWriter`.
The default context is a cancel context and after finishing the request, the cancel function is triggered and you can watch for `Done` channel in `Context()` if you need to know the request is delivered and its time to give up.

```go
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
```
See [this article](https://blog.golang.org/context) about context.

There is some example, not fully tested middleware in `middlewares` sub package, but just use them as example, anything may change there.
