package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/faruqisan/daily/pkg/cache"
	"github.com/faruqisan/daily/src/secret"
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

// FetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func (e *Engine) FetchProfile(ctx context.Context, code string) (string, error) {

	tok, err := e.oauthConf.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("exchange token: %s: %w", code, err)
	}

	peopleService, err := people.NewService(ctx, option.WithTokenSource(e.oauthConf.TokenSource(ctx, tok)))
	if err != nil {
		return "", fmt.Errorf("new people service: %w", err)
	}

	people, err := peopleService.People.Get("people/me").PersonFields("emailAddresses,photos").Do()
	if err != nil {
		return "", fmt.Errorf("get people: %w", err)
	}

	return people.EmailAddresses[0].Value, nil
}
