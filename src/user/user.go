package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/faruqisan/daily/pkg/cache"
	"github.com/jmoiron/sqlx"
)

var (
	sessionLifeTime = time.Hour * 24 // 1 day
)

type (
	// Data struct define user data
	Data struct {
		ID          int64     `db:"id" json:"id"`
		Email       string    `db:"email" json:"email"`
		LoginMethod string    `db:"login_method" json:"login_method"`
		CreatedAt   time.Time `db:"created_at" json:"created_at"`
	}

	// Engine struct define user engine to access data
	Engine struct {
		db    *sqlx.DB
		cache cache.Engine
	}
)

// New function return engine with setuped db
func New(db *sqlx.DB, cache cache.Engine) Engine {
	return Engine{
		db:    db,
		cache: cache,
	}
}

// Register function will create a new user on db
func (e Engine) Register(ctx context.Context, email, loginMethod string) error {

	// check if user already registered
	u, err := e.Login(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if u.ID != 0 {
		return errors.New("user already registered")
	}

	q := `
	INSERT INTO users
		(email, login_method)
	VALUES ($1, $2)
	`

	_, err = e.db.ExecContext(ctx, q, email, loginMethod)
	return err
}

// Login function will look into user to db
func (e Engine) Login(ctx context.Context, email string) (Data, error) {

	var (
		user Data
		err  error
	)

	q := `
	SELECT id, email, login_method, created_at FROM users WHERE email = $1
	`

	err = e.db.GetContext(ctx, &user, q, email)
	if err != nil {
		return user, err
	}

	return user, err
}

// SetSessionToCache store user session into cache
func (e Engine) SetSessionToCache(sesionKey string, user Data) error {

	var (
		byteUser []byte
		err      error
	)

	byteUser, err = json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = e.cache.SetNX(sesionKey, string(byteUser), sessionLifeTime).Result()

	return err
}

// CheckSession function check if session is exist, returning session's user
func (e *Engine) CheckSession(sessionKey string) (Data, error) {
	var (
		jsonUser string
		user     Data
		err      error
	)

	jsonUser, err = e.cache.Get(sessionKey).Result()
	if err != nil {
		return user, err
	}

	err = json.Unmarshal([]byte(jsonUser), &user)

	return user, err
}
