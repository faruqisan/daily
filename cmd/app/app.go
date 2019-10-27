package main

import (
	"log"

	"github.com/faruqisan/daily/pkg/cache"
	"github.com/faruqisan/daily/src/config"
	"github.com/faruqisan/daily/src/daily"
	"github.com/faruqisan/daily/src/server/api"
	"github.com/faruqisan/daily/src/user"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sqlx.Connect("postgres", cfg.DatabaseConfig.DSN)
	if err != nil {
		log.Println("fail to connect postgresql db : ", err)
	}

	cache := cache.New(cfg.RedisConfig.Host)

	daily := daily.New(db)
	user := user.New(db, cache)

	apiEngine := api.New(user, daily)

	log.Fatal(apiEngine.ServeHTPP())

}
