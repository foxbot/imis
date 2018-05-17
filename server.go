package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

const (
	minExpires = 1
	maxExpires = 60000
)

var (
	defaultExpires      = 15 * time.Second
	expireAfterRangeErr = fmt.Sprintf("X-Delete-After must be within (%d, %d)", minExpires, maxExpires)
)

// Server contains the server's state
type Server struct {
	Cache map[string]*[]byte
	lock  sync.RWMutex
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
	server.lock.Lock()
	server.Cache[id] = &data
	server.lock.Unlock()
	w.WriteHeader(http.StatusNoContent)

	expire := defaultExpires
	if e := r.Header.Get("X-Delete-After"); e != "" {
		d, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			http.Error(w, "X-Delete-After should be an integer", http.StatusBadRequest)
			return
		}
		if d < minExpires || d > maxExpires {
			http.Error(w, expireAfterRangeErr, http.StatusBadRequest)
			return
		}
		expire = time.Duration(d) * time.Millisecond
	}

	go func() {
		time.Sleep(expire)

		server.lock.Lock()
		delete(server.Cache, id)
		server.lock.Unlock()
	}()
}

// Get an image from the cache
func (server *Server) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "key")
	server.lock.RLock()
	data := server.Cache[id]
	server.lock.RUnlock()
	if data == nil {
		http.Error(w, "object not found", http.StatusNotFound)
		return
	}

	_, err := w.Write(*data)
	if err != nil {
		panic(err)
	}
}
