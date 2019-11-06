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

	ctx := r.Context()

	gIDToken := r.FormValue("g_id_token")

	email, err := e.auth.ValidateGoogleCallback(ctx, gIDToken)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	// TODO: check if email already registered

	err = e.user.Register(ctx, email, auth.LoginMethodGoogle)
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

	queries := r.URL.Query()

	gIDToken := queries.Get("g_id_token")

	email, err := e.auth.ValidateGoogleCallback(ctx, gIDToken)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
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
