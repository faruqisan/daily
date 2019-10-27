package api

import (
	"github.com/faruqisan/daily/src/auth"
	"github.com/faruqisan/daily/src/daily"
	"github.com/faruqisan/daily/src/user"
)

type (
	// Engine struct ..
	Engine struct {
		auth  *auth.Engine
		user  user.Engine
		daily daily.Engine
	}
)

// New function return setuped API
func New(user user.Engine, daily daily.Engine, auth *auth.Engine) Engine {
	return Engine{
		user:  user,
		daily: daily,
		auth:  auth,
	}
}
