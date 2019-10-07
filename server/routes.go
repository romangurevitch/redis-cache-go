package server

// Routes handlers
func (s *server) routes() {
	s.router.Handle("/contact", s.log(s.postContactHandler()))
	s.router.Handle("/contact/", s.log(s.getContactHandler()))
}
