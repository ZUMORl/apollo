package db

import (
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type (
	Device struct {
		Name  string `json:"name"`
		Model string `json:"model"`
	}

	Devices interface {
		Add(*Device) (string, error)
		Read(string) (Device, error)
		Update(string, *Device) error
		Delete(string) error
		List() (map[string]Device, error)
	}

	devicesManager struct {
		db DataBase
	}
)

func NewDevices(db DataBase) Devices {
	return &devicesManager{db: db}
}

func (dm *devicesManager) Add(dvc *Device) (string, error) {
	json, err := json.Marshal(dvc)
	if err != nil {
		return "", err
	}
	var key = fmt.Sprint(uuid.NewV4())
	err = dm.db.cli.Set(fmt.Sprintf("devices:%v", key), json, 0).Err()
	return key, err
}

func (dm *devicesManager) Read(key string) (dvc Device, err error) {
	val, err := dm.db.cli.Get(fmt.Sprintf("devices:%v", key)).Result()
	if err != nil {
		return Device{}, err
	}

	err = json.Unmarshal([]byte(val), &dvc)
	return dvc, err
}

func (dm *devicesManager) Update(key string, dvc *Device) error {
	json, err := json.Marshal(dvc)
	if err != nil {
		return err
	}
	err = dm.db.cli.Set(fmt.Sprintf("devices:%v", key), json, 0).Err()
	return err
}

func (dm *devicesManager) Delete(key string) error {
	return dm.db.cli.Del("devices:" + key).Err()
}

func (dm *devicesManager) List() (map[string]Device, error) {
	var keys, _, err = dm.db.cli.Scan(0, "devices:*", 0).Result()
	if err != nil || len(keys) == 0 {
		return map[string]Device{}, err
	}
	arr, err := dm.db.cli.MGet(keys...).Result()
	if err != nil {
		return map[string]Device{}, err
	}

	var ret = map[string]Device{}

	for i, elem := range arr {
		var dvc Device
		err = json.Unmarshal([]byte(elem.(string)), &dvc)
		if err != nil {
			return ret, err
		}

		ret[strings.Split(keys[i], ":")[1]] = dvc
	}
	return ret, err
}
