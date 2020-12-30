package api

func (s *Server) routes() {
	s.router.HandleFunc("/etymology", s.handleEtymology())
}
