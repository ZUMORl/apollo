package db

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"
)

const (
	sns_type  = "light"
	sns_model = "model1"
	dvc_key   = "1337"
)

func TestCreteSensor(t *testing.T) {
	var json = fmt.Sprintf("{\"Type\":\"%v\",\"Model\":\"%v\"}", sns_type, sns_model)

	ping(t)
	var sensors = NewSensors(Db)

	var key, err = sensors.Add(&Sensor{Type: sns_type, Model: sns_model}, dvc_key)
	if err != nil {
		t.Fatalf("Addition failed : %v", err)
	}

	var tst_key = "sensors:" + key + ":device:" + dvc_key
	var val = get(t, tst_key)
	assertEqual(t, val, json, "")

	delete(t, tst_key)
}

func TestReadSensor(t *testing.T) {
	var tst_id = "384"
	var json = fmt.Sprintf("{\"Type\":\"%v\",\"Model\":\"%v\"}", sns_type, sns_model)
	var tst_key = "sensors:" + tst_id + ":device:" + dvc_key

	ping(t)
	var sensors = NewSensors(Db)

	set(t, tst_key, json)

	var sns, err = sensors.Read(tst_id, dvc_key)
	if err != nil {
		t.Fatalf("Failed read : %v", err)
	}

	assertEqual(t, sns.Type, sns_type, "")
	assertEqual(t, sns.Model, sns_model, "")

	delete(t, tst_key)
}

func TestUpdateSensor(t *testing.T) {
	var tst_id = "384"
	var sns_type_ch = "Changed Type"
	var sns_model_ch = "Changed Model"
	var json = fmt.Sprintf("{\"Type\":\"%v\",\"Model\":\"%v\"}", sns_type, sns_model)
	var exp_json = fmt.Sprintf("{\"Type\":\"%v\",\"Model\":\"%v\"}", sns_type_ch, sns_model_ch)
	var tst_key = "sensors:" + tst_id + ":device:" + dvc_key

	ping(t)
	var sensors = NewSensors(Db)

	set(t, tst_key, json)

	if err := sensors.Update(tst_id, dvc_key, &Sensor{Type: sns_type_ch, Model: sns_model_ch}); err != nil {
		t.Fatalf("Failed update : %v", err)
	}

	var val = get(t, tst_key)
	assertEqual(t, val, exp_json, "")

	delete(t, tst_key)
}

func TestDeleteSensor(t *testing.T) {
	var tst_id = "384"
	var json = fmt.Sprintf("{\"Type\":\"%v\",\"Model\":\"%v\"}", sns_type, sns_model)
	var tst_key = "sensors:" + tst_id + ":device:" + dvc_key

	ping(t)
	set(t, tst_key, json)

	var sensors = NewSensors(Db)
	if err := sensors.Delete(tst_id, dvc_key); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	if exists(t, tst_key) {
		t.Fatal("Key still exists")
	}
}

func TestSensorComplex(t *testing.T) {
	var sns_ch = Sensor{Type: "Changed Type", Model: "Changed Model"}
	var json = fmt.Sprintf("{\"Type\":\"%v\",\"Model\":\"%v\"}", sns_type, sns_model)

	ping(t)
	var sensors = NewSensors(Db)

	var id, err = sensors.Add(&Sensor{Type: sns_type, Model: sns_model}, dvc_key)
	if err != nil {
		t.Fatalf("Addition fail : %v", err)
	}

	var tst_key = "sensors:" + id + ":device:" + dvc_key
	var val = get(t, tst_key)
	assertEqual(t, val, json, "")

	if err := sensors.Update(id, dvc_key, &sns_ch); err != nil {
		t.Fatalf("Update fail : %v", err)
	}

	sns, err := sensors.Read(id, dvc_key)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}
	assertEqual(t, sns, sns_ch, "")

	if err := sensors.Delete(id, dvc_key); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	if exists(t, tst_key) {
		t.Fatal("Key still exists")
	}
}

func TestListSensors(t *testing.T) {
	var tst_arr []Sensor
	var json_arr []string
	var pairs []interface{}
	for i := 0; i < 5; i += 1 {
		tst_arr = append(tst_arr, Sensor{
			fmt.Sprintf("type%v", i+1),
			fmt.Sprintf("model%v", i+1),
		})
		var val, _ = json.Marshal(tst_arr[i])
		json_arr = append(json_arr, string(val))

		pairs = append(pairs,
			fmt.Sprintf("sensors:test%v:device:%v", i+1, dvc_key),
			json_arr[i])
	}

	if err := Db.cli.MSet(pairs...).Err(); err != nil {
		t.Fatalf("Failed set : %v", err)
	}

	var sensors = NewSensors(Db)
	var arr, err = sensors.ListByDevice(dvc_key)
	if err != nil {
		t.Fatalf("Failed List : %v", err)
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Type < arr[j].Type
	})

	for i := range arr {
		assertEqual(t, arr[i], tst_arr[i],
			fmt.Sprintf("Elements %v are not equal\n%v != %v",
				i, arr[i], tst_arr[i]))
	}

	for i := range tst_arr {
		if err = Db.cli.Del(
			fmt.Sprintf("sensors:test%v:device:%v",
				i+1, dvc_key)).Err(); err != nil {
			t.Fatalf("Deletion fail : %v", err)
		}
	}
}
