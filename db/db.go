package db

import (
	"github.com/go-redis/redis"
)

type DataBase struct {
	cli *redis.Client
}

var (
	Db DataBase
)

func (db DataBase) Ping() error {
	_, err := db.cli.Ping().Result()
	return err
}

func init() {
	Db.cli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := Db.Ping(); err != nil {
		panic(err)
	}
}
