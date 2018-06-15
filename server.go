package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

// Blob is a byte slice
type Blob []byte

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
	Cache map[string]Blob
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
	server.Cache[id] = Blob(data)
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
	data, ok := server.Cache[id]
	server.lock.RUnlock()
	if !ok {
		http.Error(w, "object not found", http.StatusNotFound)
		return
	}

	_, err := w.Write(data)
	if err != nil {
		panic(err)
	}
}

// List all the objects from the cache in JSON
func (server *Server) List(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	i := 0

	server.lock.RLock()
	for key := range server.Cache {
		data[strconv.Itoa(i)] = key
		i++
	}
	server.lock.RUnlock()

	out, _ := json.Marshal(data)
	_, err := w.Write(out)
	if err != nil {
		panic(err)
	}
}
