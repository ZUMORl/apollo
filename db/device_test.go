package db

import (
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
