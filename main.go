// imis is an in-memory image server designed for use with NotSoBot
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// https://github.com/deansheather was here

// request timeout length
const timeout = 15 * time.Second

var host string
var token string

func init() {
	const (
		defaultHost  = "localhost:3000"
		defaultToken = "orange_juice"
	)
	flag.StringVar(&host, "host", defaultHost, "address in localhost:XXXX form")
	flag.StringVar(&token, "token", defaultToken, "authorization header for protected endpoints")
}

func main() {
	flag.Parse()

	r := chi.NewRouter()
	s := Server{
		Cache: make(map[string]*[]byte),
	}

	// builtin middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(timeout))

	// custom middlewares
	auth := AuthorizationMiddleware(token)

	r.With(auth).Post("/objects/{key}", s.Upload)
	r.Get("/objects/{key}", s.Get)

	log.Fatalln(http.ListenAndServe(host, r))
}
