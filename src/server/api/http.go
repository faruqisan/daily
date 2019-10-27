package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/faruqisan/daily/src/auth"
	"github.com/faruqisan/daily/src/server/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
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
	r.Get("/login/google-callback", e.handleGoogleLoginCallback)
	r.Get("/register", e.handleRegister)

	// require authorization endpoints
	r.Group(func(r chi.Router) {
		r.Use(e.Authorization())
		r.Get("/ping-auth", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong-auth"))
		})
	})

	// r.Route("/job", func(r chi.Router) {
	// 	r.Post("/http", s.HandlePushHTTPJob)
	// })

	return r
}

func (e Engine) handleRegister(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	redirectURL, err := e.auth.GenerateGoogleURL(auth.ActionRegister)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	d := struct {
		RedirectURL string `json:"redirect_url"`
	}{
		RedirectURL: redirectURL,
	}

	resp.Data = d
}

func (e Engine) handleLogin(w http.ResponseWriter, r *http.Request) {

	resp := response.Response{}
	defer resp.Render(w, r)

	redirectURL, err := e.auth.GenerateGoogleURL(auth.ActionLogin)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	d := struct {
		RedirectURL string `json:"redirect_url"`
	}{
		RedirectURL: redirectURL,
	}

	resp.Data = d
}

func (e Engine) handleGoogleLoginCallback(w http.ResponseWriter, r *http.Request) {

	resp := response.Response{}
	defer resp.Render(w, r)

	ctx := r.Context()

	state := r.FormValue("state")
	code := r.FormValue("code")

	action, email, err := e.auth.ValidateGoogleCallback(ctx, state, code)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	if action == auth.ActionRegister {
		err = e.user.Register(ctx, email, auth.LoginMethodGoogle)
		if err != nil {
			resp.SetError(err, http.StatusInternalServerError)
			return
		}
	}

	// check if user registered
	user, err := e.user.Login(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows || user.Email == "" {
		resp.SetError(errors.New("user not registered, please signup"), http.StatusUnauthorized)
		return
	}

	// set session
	session := uuid.New().String()
	err = e.user.SetSessionToCache(session, user)

	d := struct {
		Success bool
		Session string
	}{
		Success: true,
		Session: session,
	}

	resp.Data = d

}
