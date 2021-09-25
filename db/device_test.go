package db

import (
	"fmt"
	"testing"
)

const (
	dvcName  = "device1"
	dvcModel = "model1"
)

func TestDeviceComplex(t *testing.T) {
	var dvcStart = Device{Name: dvcName, Model: dvcModel}
	var dvcChanged = Device{Name: "Changed Name", Model: "Changed Model"}

	var devices = NewDevices(Db)

	var id, err = devices.Add(&dvcStart)
	if err != nil {
		t.Fatalf("Addition fail : %v", err)
	}

	dvc, err := devices.Read(id)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}
	assertEqual(t, dvc, dvcStart, "")

	if err := devices.Update(id, &dvcChanged); err != nil {
		t.Fatalf("Update fail : %v", err)
	}

	dvc, err = devices.Read(id)
	if err != nil {
		t.Fatalf("Read fail : %v", err)
	}
	assertEqual(t, dvc, dvcChanged, "")

	if err := devices.Delete(id); err != nil {
		t.Fatalf("Deletion fail : %v", err)
	}

	if _, err = devices.Read(id); err == nil {
		t.Fatal("Key still exists")
	}
}

func TestListDevices(t *testing.T) {
	var num = 5
	var tstArr = make([]Device, num)
	var ids = make([]string, num)
	var devices = NewDevices(Db)
	for i := 0; i < num; i++ {
		tstArr[i] = Device{
			fmt.Sprintf("name%v", i+1),
			fmt.Sprintf("model%v", i+1),
		}
		var id, err = devices.Add(&tstArr[i])
		if err != nil {
			t.Fatalf("Failed Add Device : %v", err)
		}
		ids[i] = id
	}

	var retMap, err = devices.List()
	if err != nil {
		t.Fatalf("Failed List : %v", err)
	}
	for _, id := range ids {
		if _, exists := retMap[id]; !exists {
			t.Fail()
			t.Logf("Couldn't get key %v", id)
		}
	}

	for _, id := range ids {
		if err := devices.Delete(id); err != nil {
			t.Fatalf("Deletion fail : %v", err)
		}
	}
}
