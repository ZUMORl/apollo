package db

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

const (
	tst_name  = "name1"
	tst_model = "model1"
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

func TestCreateDevice(t *testing.T) {
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)

	ping(t)
	var devices = NewDevices(Db)

	var key, err = devices.Add(&Device{Name: tst_name, Model: tst_model})
	if err != nil {
		t.Fatalf("Addition failed : %v", err)
	}

	var tst_key = "devices:" + key

	var val = get(t, tst_key)
	assertEqual(t, val, json, "")

	delete(t, tst_key)
}

func TestReadDevice(t *testing.T) {
	var tst_id = "384"
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)
	var tst_key = "devices:" + tst_id

	ping(t)
	var devices = NewDevices(Db)

	set(t, tst_key, json)

	var dvc, err = devices.Read(tst_id)
	if err != nil {
		t.Fatalf("Failed read : %v", err)
	}

	assertEqual(t, dvc.Name, tst_name, "")
	assertEqual(t, dvc.Model, tst_model, "")

	delete(t, tst_key)
}

func TestUpdateDevice(t *testing.T) {
	var tst_id = "384"
	var tst_name_ch = "Changed Name"
	var tst_model_ch = "Changed Model"
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)
	var exp_json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name_ch, tst_model_ch)
	var tst_key = "devices:" + tst_id

	ping(t)
	set(t, tst_key, json)

	var devices = NewDevices(Db)
	if err := devices.Update(tst_id, &Device{Name: tst_name_ch, Model: tst_model_ch}); err != nil {
		t.Fatalf("Failed update : %v", err)
	}

	var val = get(t, tst_key)
	assertEqual(t, val, exp_json, "")

	delete(t, tst_key)
}

func TestDeleteDevice(t *testing.T) {
	var tst_id = "384"
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)
	var tst_key = "devices:" + tst_id

	ping(t)
	set(t, tst_key, json)

	var devices = NewDevices(Db)
	if err := devices.Delete(tst_id); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	_, err := Db.cli.Get(tst_key).Result()
	switch err {
	case nil:
		t.Fatal("Key still exists")
	case redis.Nil:
		break
	default:
		t.Fatalf("Failed get : %v", err)
	}
}
func TestDeviceComplex(t *testing.T) {
	// var tst_name_ch = "Changed Name"
	// var tst_model_ch = "Changed Model"
	var dvc_ch = Device{Name: "Changed Name", Model: "Changed Model"}
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name, tst_model)
	// var json_ch = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", tst_name_ch, tst_model_ch)

	ping(t)
	var devices = NewDevices(Db)

	var id, err = devices.Add(&Device{Name: tst_name, Model: tst_model})
	if err != nil {
		t.Fatalf("Addition fail : %v", err)
	}

	var tst_key = "devices:" + id

	var val = get(t, tst_key)
	assertEqual(t, val, json, "")

	if devices.Update(id, &dvc_ch); err != nil {
		t.Fatalf("Update fail : %v", err)
	}

	dvc, err := devices.Read(id)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}

	assertEqual(t, dvc, dvc_ch, "")

	devices.Delete(id)

	_, err = Db.cli.Get(tst_key).Result()
	switch err {
	case nil:
		t.Fatal("Key still exists")
	case redis.Nil:
		break
	default:
		t.Fatalf("Failed get : %v", err)
	}
}
