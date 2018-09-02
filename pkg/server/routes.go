package server

import (
	"log"
	"net/http"
)

func logRequestMidleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received:", r.Method, "at:", r.URL, "from:", r.RemoteAddr)
		fn(w, r)
	}
}

func (s *serv) routes() {
	s.router.HandleFunc("/hello", logRequestMidleware(s.handleHello()))
	s.router.HandleFunc("/restaurants", logRequestMidleware(s.handleRestaurants()))
	// s.router.HandleFunc("/api/", s.handleAPI())
	// s.router.HandleFunc("/about", s.handleAbout())
	// s.router.HandleFunc("/", s.handleIndex())
	// s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex))
}
