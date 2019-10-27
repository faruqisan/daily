package user

import (
	"encoding/json"
	"context"
	"time"

	"github.com/faruqisan/daily/pkg/cache"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	sessionLifeTime = time.Hour * 24 // 1 day
)

type (
	// Data struct define user data
	Data struct {
		ID        int64     `db:"id"`
		Email     string    `db:"email"`
		Password  string    `db:"password"`
		CreatedAt time.Time `db:"created_at"`
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
		db: db,
	}
}

// Register function will create a new user on db, password will hashed inside
func (e Engine) Register(ctx context.Context, email, password string) error {
	q := `
	INSERT INTO users
		(email, password)
	VALUES ($1, $2)
	`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	_, err = e.db.ExecContext(ctx, q, email, string(hashedPassword))
	return err
}

// Login function will look into user to db
func (e Engine) Login(ctx context.Context, email, password string) (Data, error) {

	var (
		user Data
		err  error
	)

	q := `
	SELECT id, email, password, created_at FROM users WHERE email = $1
	`

	err = e.db.GetContext(ctx, &user, q, email)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// remove user password from object that will returned
	user.Password = ""

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