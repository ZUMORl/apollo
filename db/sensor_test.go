package db

import (
	"fmt"
	"sort"
	"testing"
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

	sns, err := sensors.Read(id, dvcId)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}
	assertEqual(t, sns, snsStart, "")

	if err := sensors.Update(id, dvcId, &snsChanged); err != nil {
		t.Fatalf("Update fail : %v", err)
	}

	sns, err = sensors.Read(id, dvcId)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}
	assertEqual(t, sns, snsChanged, "")

	if err := sensors.Delete(id, dvcId); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	if _, err := sensors.Read(id, dvcId); err == nil {
		t.Fatal("Key still exist")
	}
}

func TestListSensors(t *testing.T) {
	var tstArr []Sensor
	var sensors = NewSensors(Db)
	var ids []string
	var num = 5
	for i := 0; i < num; i += 1 {
		tstArr = append(tstArr, Sensor{
			fmt.Sprintf("type%v", i+1),
			fmt.Sprintf("model%v", i+1),
		})
		var id, err = sensors.Add(&tstArr[i], dvcId)
		if err != nil {
			t.Fatalf("Failed Add Sensor : %v", err)
		}
		ids = append(ids, id)
	}

	var retIds, arr, err = sensors.ListByDevice(dvcId)
	if err != nil {
		t.Fatalf("Failed List : %v", err)
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Type < arr[j].Type
	})
	sort.Strings(retIds)
	sort.Strings(ids)

	if !compareSlices(retIds, ids) {
		t.Fatal("Initial and returned ids are not equal")
	}

	for i := range arr {
		assertEqual(t, arr[i], tstArr[i],
			fmt.Sprintf("Elements %v are not equal\n%v != %v",
				i, arr[i], tstArr[i]))
	}

	for _, id := range ids {
		if err := sensors.Delete(id, dvcId); err != nil {
			t.Fatalf("Deletion fail : %v", err)
		}
	}
}
