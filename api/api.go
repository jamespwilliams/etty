package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jamespwilliams/etty"
	"go.uber.org/zap"
)

type Server struct {
	etty   etty.Etymology
	router chi.Router
	logger *zap.Logger
}

func NewServer(logger *zap.Logger, etty etty.Etymology) Server {
	s := Server{
		etty:   etty,
		router: chi.NewRouter(),
		logger: logger,
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

		node := s.etty.Lookup(etty.Word{
			Word:     word,
			Language: lang,
		})

		if err := json.NewEncoder(w).Encode(node); err != nil {
			s.logger.Error("failed to encode response", zap.Error(err))
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
