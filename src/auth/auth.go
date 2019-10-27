package auth

import (
	"context"
	"fmt"
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
func (e *Engine) ValidateGoogleCallback(ctx context.Context, state, code string) (action string, email string, err error) {
	// check for state on cache
	cacheKey := fmt.Sprintf(googleAuthCacheKey, state)
	action, err = e.cache.Get(cacheKey).Result()
	if err != nil {
		return
	}

	tok, err := e.oauthConf.Exchange(ctx, code)
	if err != nil {
		return
	}

	people, err := e.fetchProfile(ctx, tok)
	if err != nil {
		return
	}

	email = people.EmailAddresses[0].Value

	return
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
