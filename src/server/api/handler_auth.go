package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/faruqisan/daily/src/auth"
	"github.com/faruqisan/go-response"
	"github.com/google/uuid"
)

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
