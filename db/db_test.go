package db

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func ping(t *testing.T) {
	if err := Db.Ping(); err != nil {
		t.Fatal("Can't connect to DB")
	}
}

func get(t *testing.T, key string) string {
	val, err := Db.cli.Get(key).Result()
	if err != nil {
		t.Fatalf("Failed get : %v", err)
	}
	return val
}

func set(t *testing.T, key string, val string) {
	if err := Db.cli.Set(key, val, 0).Err(); err != nil {
		t.Fatalf("Failed set : %v", err)
	}
}

func delete(t *testing.T, key string) {
	if err := Db.cli.Del(key).Err(); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		if len(msg) == 0 {
			msg = fmt.Sprintf("Unexpected result - %v, need - %v", a, b)
		}
		t.Log(msg)
		t.Fail()
	}
}

func exists(t *testing.T, key string) bool {
	_, err := Db.cli.Get(key).Result()
	switch err {
	case nil:
		return true
	case redis.Nil:
		return false
	default:
		t.Fatalf("Failed check : %v", err)
		return false
	}
}

func TestDbComplex(t *testing.T) {
	const key = "Key:123"
	const val = "{ Value : 123 }"
	if err := Db.cli.Set(key, val, 0).Err(); err != nil {
		t.Fatalf("Failed set : %v", err)
	}
	result, err := Db.cli.Get(key).Result()
	if err != nil {
		t.Fatalf("Failed get : %v", err)
	}
	if result != val {
		t.Logf("Get value not equal to set : \n\t%v != %v", result, val)
		t.Fail()
	}
	if err := Db.cli.Del(key).Err(); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}
}
