package db

import (
	"fmt"
	"testing"
)

const (
	dvc_name  = "device1"
	dvc_model = "model1"
)

func TestCreateDevice(t *testing.T) {
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", dvc_name, dvc_model)

	ping(t)
	var devices = NewDevices(Db)

	var key, err = devices.Add(&Device{Name: dvc_name, Model: dvc_model})
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
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", dvc_name, dvc_model)
	var tst_key = "devices:" + tst_id

	ping(t)
	var devices = NewDevices(Db)

	set(t, tst_key, json)

	var dvc, err = devices.Read(tst_id)
	if err != nil {
		t.Fatalf("Failed read : %v", err)
	}

	assertEqual(t, dvc.Name, dvc_name, "")
	assertEqual(t, dvc.Model, dvc_model, "")

	delete(t, tst_key)
}

func TestUpdateDevice(t *testing.T) {
	var tst_id = "384"
	var dvc_name_ch = "Changed Name"
	var dvc_model_ch = "Changed Model"
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", dvc_name, dvc_model)
	var exp_json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", dvc_name_ch, dvc_model_ch)
	var tst_key = "devices:" + tst_id

	ping(t)
	set(t, tst_key, json)

	var devices = NewDevices(Db)
	if err := devices.Update(tst_id, &Device{Name: dvc_name_ch, Model: dvc_model_ch}); err != nil {
		t.Fatalf("Failed update : %v", err)
	}

	var val = get(t, tst_key)
	assertEqual(t, val, exp_json, "")

	delete(t, tst_key)
}

func TestDeleteDevice(t *testing.T) {
	var tst_id = "384"
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", dvc_name, dvc_model)
	var tst_key = "devices:" + tst_id

	ping(t)
	set(t, tst_key, json)

	var devices = NewDevices(Db)
	if err := devices.Delete(tst_id); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	if exists(t, tst_key) {
		t.Fatal("Key still exists")
	}
}
func TestDeviceComplex(t *testing.T) {
	var dvc_ch = Device{Name: "Changed Name", Model: "Changed Model"}
	var json = fmt.Sprintf("{\"Name\":\"%v\",\"Model\":\"%v\"}", dvc_name, dvc_model)

	ping(t)
	var devices = NewDevices(Db)

	var id, err = devices.Add(&Device{Name: dvc_name, Model: dvc_model})
	if err != nil {
		t.Fatalf("Addition fail : %v", err)
	}

	var tst_key = "devices:" + id

	var val = get(t, tst_key)
	assertEqual(t, val, json, "")

	if err := devices.Update(id, &dvc_ch); err != nil {
		t.Fatalf("Update fail : %v", err)
	}

	dvc, err := devices.Read(id)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}

	assertEqual(t, dvc, dvc_ch, "")

	if err := devices.Delete(id); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	if exists(t, tst_key) {
		t.Fatal("Key still exists")
	}
}
