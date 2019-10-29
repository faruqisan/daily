package api

import (
	"errors"
	"net/http"

	"github.com/faruqisan/daily/pkg/session"
	"github.com/faruqisan/go-response"
)

var (
	errUnauthorized   = errors.New("session expired")
	errSessionMissing = errors.New("session missing from header")
)

// Authorization function will check given session key if logedin or not
func (e Engine) Authorization() func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			resp := response.Response{}
			sessionKey := r.Header.Get("session")

			if sessionKey == "" {
				resp.SetError(errSessionMissing, http.StatusBadRequest)
				resp.Render(w, r)
				return
			}

			userData, err := e.user.CheckSession(sessionKey)
			if err != nil {
				resp.SetError(err, http.StatusBadRequest)
				resp.Render(w, r)
				return
			}

			// check if not exist
			if userData.ID == 0 {
				resp.SetError(errUnauthorized, http.StatusUnauthorized)
				resp.Render(w, r)
				return
			}

			contextWithUID := session.SetUIDToCtx(r.Context(), userData.ID)

			next.ServeHTTP(w, r.WithContext(contextWithUID))
		}

		return http.HandlerFunc(fn)
	}

}
