package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/faruqisan/go-response"
	"github.com/google/uuid"
)

func (e Engine) handleRegister(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	ctx := r.Context()

	// check for code that received from google auth callback
	code := r.FormValue("code")
	if code == "" {
		resp.SetError(errors.New("code is missing"))
		return
	}

	// get user email from code
	email, err := e.auth.FetchProfile(ctx, code)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	err = e.user.Register(ctx, email, "google")
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	resp.SetSuccess()
}

func (e Engine) handleLogin(w http.ResponseWriter, r *http.Request) {

	resp := response.Response{}
	defer resp.Render(w, r)

	ctx := r.Context()

	// check for code that received from google auth callback
	code := r.FormValue("code")
	if code == "" {
		resp.SetError(errors.New("code is missing"))
		return
	}

	// get user email from code
	email, err := e.auth.FetchProfile(ctx, code)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	// check on db
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
