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
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)
	var tst_key = fmt.Sprintf("devices:%v", tst_id)

	var cli = connect()
	var dvc = NewDevice(tst_id, tst_name, tst_model)
	var err = dvc.Create(cli)
	if err != nil {
		t.Fatalf("Failed create : %v", err)
	}

	val, err := cli.Get(tst_key).Result()
	if err != nil {
		t.Fatalf("Failed get : %v", err)
	}

	if val != json {
		t.Fatalf("Unexpected result - %v, need - %v", val, json)
	}

	if err := cli.Del(tst_key).Err(); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}
}

func TestReadDevice(t *testing.T) {
	var tst_id = 384
	var tst_name = "name1"
	var tst_model = "model1"
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)
	var tst_key = fmt.Sprintf("devices:%v", tst_id)

	var cli = connect()
	var err = cli.Set(tst_key, json, 0).Err()
	if err != nil {
		t.Fatalf("Failed set : %v", err)
	}
	var dvc = NewDevice(tst_id)
	err = dvc.Read(cli)
	if err != nil {
		t.Fatalf("Failed read : %v", err)
	}
	if dvc.Name != tst_name && dvc.Model != tst_model {
		t.Logf("Unexpected result : %v, %v \n\texpected : %v, %v",
			dvc.Name, dvc.Model, tst_name, tst_model)
		t.Fail()
	}

	if err := cli.Del(tst_key).Err(); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}
}
