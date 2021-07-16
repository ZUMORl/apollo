package db

import (
	"testing"
)

func TestPing(t *testing.T) {
	if err := Db.Ping(); err != nil {
		t.Fatal("Can't connect to DB")
	}
}

// func TestGet(t *testing.T) {
// 	val, err := Db.cli.Get(key).Result()
// 	if err != nil {
// 		t.Fatalf("Failed get : %v", err)
// 	}
// }

// func TestSet(t *testing.T) {
// 	if err := Db.cli.Set(key, val, 0).Err(); err != nil {
// 		t.Fatalf("Failed set : %v", err)
// 	}
// }

// func TestDelete(t *testing.T) {
// 	if err := Db.cli.Del(key).Err(); err != nil {
// 		t.Fatalf("Deletion fail : %v", err)
// 	}
// }

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
