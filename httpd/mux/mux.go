// Package mux provides ...
package mux

import "net/http"

// handler
type Mux struct {
	serveMux *http.ServeMux
}

// global variable
var Muxd = NewMux()

// new mux
func NewMux() *Mux {
	return &Mux{
		serveMux: http.NewServeMux(),
	}
}

// implement handler
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	m.serveMux.ServeHTTP(w, r)
}

// handle func
func (m *Mux) Handle(pattern string, handler HandleFunc) {
	m.serveMux.Handle(pattern, handler)
}

// handle func
func Handle(pattern string, handler HandleFunc) {
	Muxd.Handle(pattern, handler)
}
