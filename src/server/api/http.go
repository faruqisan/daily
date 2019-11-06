package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// ServeHTPP function will run http server
// run this on main app
func (e Engine) ServeHTPP() error {
	return http.ListenAndServe(":8080", e.handler())
}

func (e Engine) handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/login", e.handleLogin)
	r.Get("/register", e.handleRegister)

	// require authorization endpoints
	r.Group(func(r chi.Router) {
		r.Use(e.Authorization())
		r.Get("/ping-auth", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong-auth"))
		})

		r.Route("/reports", func(r chi.Router) {
			r.Post("/", e.handleCreateReports)
			r.Get("/search", e.handleGetUserReports)
		})

	})

	return r
}
