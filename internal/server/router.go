package server

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("/_healthz", Healthz).Methods("GET").Name("Healthz")
}
