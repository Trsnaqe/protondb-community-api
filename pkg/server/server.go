package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/protondb-api/pkg/api"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	server := &Server{
		router: mux.NewRouter(),
	}

	api.SetupRoutes(server.router)

	return server
}

func (s *Server) Run(addr string) {
	log.Printf("Server started at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.router))
}
