package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/faruqisan/daily/pkg/cache"
	"github.com/faruqisan/daily/src/secret"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type (
	// Engine ..
	Engine struct {
		cache     cache.Engine
		oauthConf *oauth2.Config
		secret    secret.Secret
	}

	tokenInfo struct {
		Iss string `json:"iss"`
		Sub string `json:"sub"`
		Azp string `json:"azp"`
		Aud string `json:"aud"`
		Iat string `json:"iat"`
		Exp string `json:"exp"`

		Email         string `json:"email,omitempty"`
		EmailVerified string `json:"email_verified,omitempty"`
		Name          string `json:"name,omitempty"`
		Picture       string `json:"picture,omitempty"`
		GivenName     string `json:"given_name,omitempty"`
		FamilyName    string `json:"family_name,omitempty"`
		Locale        string `json:"locale,omitempty"`

		ErrorDescription string `json:"error_description,omitempty"`
	}
)

const (
	// gauth:uuid , value should be action login/register
	googleAuthCacheKey = "gauth:%s"
	// ActionLogin const define login action
	ActionLogin = "login"
	// ActionRegister const define register action
	ActionRegister = "register"
	// LoginMethodGoogle const define login method using google
	LoginMethodGoogle = "google"

	googleTokenInfoURL = "https://www.googleapis.com/oauth2/v3/tokeninfo"
)

var (
	sessionExpire = 10 * time.Minute
)

// New function
func New(sec secret.Secret, cache cache.Engine) *Engine {
	conf := &oauth2.Config{
		ClientID:     sec.GoogleOAuth.ClientID,
		ClientSecret: sec.GoogleOAuth.ClientSecret,
		RedirectURL:  sec.GoogleOAuth.RedirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return &Engine{
		cache:     cache,
		oauthConf: conf,
		secret:    sec,
	}
}

// GenerateGoogleURL will generate redirect url to google auth
// and save cache based on action given
func (e *Engine) GenerateGoogleURL(action string) (string, error) {

	sessionID := uuid.New().String()
	cacheKey := fmt.Sprintf(googleAuthCacheKey, sessionID)

	err := e.cache.Set(cacheKey, action, sessionExpire).Err()
	if err != nil {
		return "", err
	}

	return e.oauthConf.AuthCodeURL(sessionID), nil
}

// ValidateGoogleCallback ..
func (e *Engine) ValidateGoogleCallback(ctx context.Context, idToken string) (string, error) {

	resp, err := http.PostForm(googleTokenInfoURL, url.Values{"id_token": {idToken}})
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status not OK")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	tokenInfo := tokenInfo{}

	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return "", err
	}

	if tokenInfo.ErrorDescription != "" {
		return "", errors.New(tokenInfo.ErrorDescription)
	}

	if tokenInfo.Aud != e.secret.GoogleOAuth.ClientID {
		log.Println(tokenInfo.Aud)
		log.Println(e.secret.GoogleOAuth.ClientID)
		return "", errors.New("client id not match")
	}

	return tokenInfo.Email, nil
}

// fetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func (e *Engine) fetchProfile(ctx context.Context, tok *oauth2.Token) (*people.Person, error) {
	peopleService, err := people.NewService(ctx, option.WithTokenSource(e.oauthConf.TokenSource(ctx, tok)))
	if err != nil {
		return nil, err
	}
	return peopleService.People.Get("people/me").PersonFields("emailAddresses,photos").Do()
}
