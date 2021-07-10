package entities

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func connect() *redis.Client {
	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func TestPing(t *testing.T) {
	var cli = connect()
	_, err := cli.Ping().Result()
	if err != nil {
		t.Fatalf("No connection: %v", err)
	}
}

func TestCreateDevice(t *testing.T) {
	var tst_id = 384
	var tst_name = "name1"
	var tst_model = "model1"
	var exp = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)
	var tst_key = fmt.Sprintf("devices:%v", tst_id)

	var cli = connect()
	var dvc = NewDevice(tst_id, tst_name, tst_model)
	dvc.Create(cli)

	val, err := cli.Get(tst_key).Result()
	if err != nil {
		t.Fatalf("Failed get : %v", err)
	}
	if val != exp {
		t.Fatalf("Unexpected result - %v, need - %v", val, exp)
	}
	if err := cli.Del(tst_key).Err(); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}
}
