package server

import (
	"cli_example/app"
	"net/http"
)

// Server is the interface for the server
type Server interface {
	Start() error
}

// server is the implementation of the server
type server struct {
	address string
	app     app.Application
}

// NewServer creates a new server
func NewServer(address string, app app.Application) Server {
	return &server{
		address: address,
		app:     app,
	}
}

// Start starts the server and waits for it to stop
func (s *server) Start() error {
	http.HandleFunc("/hello", s.hello)
	http.HandleFunc("/version", s.Version)

	return http.ListenAndServe(s.address, nil)
}

// hello is the handler for the hello route
func (s *server) hello(r http.ResponseWriter, req *http.Request) {
	r.Write([]byte("Hello World"))
}

// Version is the handler for the version route
func (s *server) Version(r http.ResponseWriter, req *http.Request) {
	r.Write([]byte(s.app.Version()))
}
