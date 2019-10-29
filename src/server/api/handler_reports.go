package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/faruqisan/daily/pkg/session"
	"github.com/faruqisan/go-response"
)

const (
	dailyTimeLayout = "02-01-2006"
)

func (e Engine) handleCreateReports(w http.ResponseWriter, r *http.Request) {

	resp := response.Response{}
	defer resp.Render(w, r)

	ctx := r.Context()

	uid, err := session.GetUIDFromCTX(ctx)
	if err != nil {
		resp.SetError(err, http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	detail := r.FormValue("detail")

	err = e.daily.Create(ctx, uid, title, detail)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	resp.SetSuccess()

}

func (e Engine) handleGetUserReports(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	ctx := r.Context()

	uid, err := session.GetUIDFromCTX(ctx)
	if err != nil {
		resp.SetError(err, http.StatusUnauthorized)
		return
	}

	urlQueries := r.URL.Query()
	rawTStart := urlQueries["time_start"][0]
	rawTEnd := urlQueries["time_end"][0]
	if rawTStart == "" {
		resp.SetError(errors.New("time can't blank"), http.StatusBadRequest)
		return
	}

	tStart, err := time.Parse(dailyTimeLayout, rawTStart)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	tEnd, err := time.Parse(dailyTimeLayout, rawTEnd)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	tEndEod := time.Date(tEnd.Year(), tEnd.Month(), tEnd.Day()+1, 0, 0, 0, -1, tEnd.Location())
	reports, err := e.daily.GetUserReports(ctx, uid, tStart, tEndEod)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		return
	}

	resp.Data = reports
}
