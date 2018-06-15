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

const (
	defaultHost  = "0.0.0.0:3000"
	defaultToken = "orange_juice"
)

var (
	host  = defaultHost
	token = defaultToken
)

func init() {
	flag.StringVar(&host, "host", defaultHost, "address in 0.0.0.0:0 form")
	flag.StringVar(&token, "token", defaultToken, "authorization header for protected endpoints")
}

func main() {
	flag.Parse()

	handler := buildRouter()

	log.Fatalln(http.ListenAndServe(host, handler))
}

func buildRouter() http.Handler {
	r := chi.NewRouter()
	s := Server{
		Cache: make(map[string]Blob),
	}

	// builtin middlewares
	r.Use(middleware.GetHead)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(timeout))

	// custom middlewares
	auth := AuthorizationMiddleware(token)

	r.With(auth).Post("/objects/{key}", s.Upload)
	r.With(auth).Get("/objects", s.List)
	r.Get("/objects/{key}", s.Get)

	return r
}
