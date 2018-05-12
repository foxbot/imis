package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
)

// Server contains the server's state
type Server struct {
	Cache map[string]*[]byte
}

// Upload an image to the cache
func (server *Server) Upload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "key")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		http.Error(w, "payload was empty", http.StatusBadRequest)
		return
	}
	server.Cache[id] = &data
	w.WriteHeader(http.StatusNoContent)
}

// Get an image from the cache
func (server *Server) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "key")
	data := server.Cache[id]
	if data == nil {
		http.Error(w, "object not found", http.StatusNotFound)
		return
	}

	_, err := w.Write(*data)
	if err != nil {
		panic(err)
	}

	delete(server.Cache, id)
}
