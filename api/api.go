package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jamespwilliams/etymology"
)

type Server struct {
	ety    etymology.Etymology
	router chi.Router
}

func NewServer(ety etymology.Etymology) Server {
	s := Server{
		ety:    ety,
		router: chi.NewRouter(),
	}
	s.routes()
	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleEtymology() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world\n")
	}
}
