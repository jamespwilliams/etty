package api

import (
	"encoding/json"
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
		if r.Method != http.MethodGet {
			httpError(w, http.StatusBadRequest)
			return
		}

		word := r.URL.Query().Get("word")
		if word == "" {
			httpError(w, http.StatusBadRequest, "missing word parameter")
			return
		}

		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = "en"
		}

		node := s.ety.Lookup(etymology.Word{
			Word:     word,
			Language: lang,
		})

		if err := json.NewEncoder(w).Encode(node); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func httpError(w http.ResponseWriter, code int, message ...string) {
	if len(message) > 0 {
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(code), message[0]), code)
		return
	}

	http.Error(w, http.StatusText(code), code)
}
