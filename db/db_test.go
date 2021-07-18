package db

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		if len(msg) == 0 {
			msg = fmt.Sprintf("Unexpected result - %v, need - %v", a, b)
		}
		t.Log(msg)
		t.Fail()
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

	assertEqual(t, result, val, "")

	if err := Db.cli.Del(key).Err(); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}
}
