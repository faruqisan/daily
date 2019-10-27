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
	r.Post("/register", e.handleRegister)

	// r.Route("/job", func(r chi.Router) {
	// 	r.Post("/http", s.HandlePushHTTPJob)
	// })

	return r
}

func (e Engine) handleRegister(w http.ResponseWriter, r *http.Request) {

}

func (e Engine) handleLogin(w http.ResponseWriter, r *http.Request) {}
