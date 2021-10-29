package db

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

const (
	sensorType  = "light"
	sensorModel = "model1"
	dvcId       = "1337"
)

func compareSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSensorComplex(t *testing.T) {
	var snsStart = Sensor{Type: sensorType, Model: sensorModel}
	var snsChanged = Sensor{Type: "Changed Type", Model: "Changed Model"}

	var sensors = NewSensors(Db)

	var id, err = sensors.Add(&Sensor{Type: sensorType, Model: sensorModel}, dvcId)
	if err != nil {
		t.Fatalf("Addition fail : %v", err)
	}

	sns, err := sensors.Read(id)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}
	assertEqual(t, sns, snsStart, "")

	if err := sensors.Update(id, &snsChanged); err != nil {
		t.Fatalf("Update fail : %v", err)
	}

	sns, err = sensors.Read(id)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}
	assertEqual(t, sns, snsChanged, "")

	if err := sensors.Delete(id); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	if _, err := sensors.Read(id); err == nil {
		t.Fatal("Key still exist")
	}
}

func TestListSensors(t *testing.T) {
	var num = 5
	var tstArr = make([]Sensor, num)
	var ids = make([]string, num)
	var sensors = NewSensors(Db)
	for i := 0; i < num; i++ {
		tstArr[i] = Sensor{
			fmt.Sprintf("type%v", i+1),
			fmt.Sprintf("model%v", i+1),
		}
		var id, err = sensors.Add(&tstArr[i], dvcId)
		if err != nil {
			t.Fatalf("Failed Add Sensor : %v", err)
		}
		ids[i] = id
	}

	var retMap, err = sensors.ListByDevice(dvcId)
	if err != nil {
		t.Fatalf("Failed List : %v", err)
	}

	var retIds = make([]string, len(retMap))
	var i = 0
	for id := range retMap {
		retIds[i] = id
		i++
	}
	sort.Strings(retIds)
	sort.Strings(ids)

	if !compareSlices(retIds, ids) {
		t.Fatal("Initial and returned ids are not equal")
	}

	for _, id := range ids {
		if err := sensors.Delete(id); err != nil {
			t.Fatalf("Deletion fail : %v", err)
		}
	}
}

func TestValuesComplex(t *testing.T) {
	var sensors = NewSensors(Db)
	var snsTest = &Sensor{Type: "Test", Model: "Test-v1"}
	var id, err = sensors.Add(snsTest, dvcId)
	if err != nil {
		t.Fatalf("Add sensor fail : %v", err)
	}
	var timeformat = "15:04:05 02.01.2006"

	var testValues = []Value{
		Value{Val: "test1", Timestamp: time.Now().Format(timeformat)},
		Value{Val: "test2", Timestamp: (time.Now().Add(time.Second)).Format(timeformat)},
		Value{Val: "test2", Timestamp: (time.Now().Add(time.Second * 5)).Format(timeformat)},
	}
	for _, testVal := range testValues {
		if err := sensors.AddValue(id, &testVal); err != nil {
			t.Fatalf("Add fail : %v", err)
		}
	}

	retVals, err := sensors.GetValues(id, 0, -1)
	if err != nil {
		t.Fatalf("Get fail : %v", err)
	}

	testValues[0], testValues[2] = testValues[2], testValues[0]
	if reflect.DeepEqual(retVals, testValues) != true {
		t.Fail()
		t.Logf("Got wrong values")
	}

	if err = sensors.Delete(id); err != nil {
		t.Fatalf("Clean up fail : %v", err)
	}

	if err = Db.cli.Del("values:sensors:" + id + ":device:" + dvcId).Err(); err != nil {
		t.Fatalf("Clean up fail : %v", err)
	}
}
