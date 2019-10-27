package api

import (
	"github.com/faruqisan/daily/src/daily"
	"github.com/faruqisan/daily/src/user"
)

type (
	// Engine struct ..
	Engine struct {
		user  user.Engine
		daily daily.Engine
	}
)

// New function return setuped API
func New(user user.Engine, daily daily.Engine) Engine {
	return Engine{
		user:  user,
		daily: daily,
	}
}
