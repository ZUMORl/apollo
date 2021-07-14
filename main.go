package main

import (
	_ "github.com/apollo/entities"

	"fmt"
	"os"

	"github.com/go-redis/redis"
)

func checkConnection(cli *redis.Client) error {
	_, err := cli.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Connected successfuly")
	return nil
}

func main() {
	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := checkConnection(client); err != nil {
		os.Exit(1)
	}
}
