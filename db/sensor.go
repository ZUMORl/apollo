package db

import (
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type (
	Sensor struct {
		Type  string
		Model string
	}

	Sensors interface {
		Add(*Sensor, string) (string, error)
		Read(string, string) (Sensor, error)
		Update(string, string, *Sensor) error
		Delete(string, string) error
		ListByDevice(string) ([]string, []Sensor, error)
	}

	sensorManager struct {
		db DataBase
	}
)

func NewSensors(db DataBase) Sensors {
	return &sensorManager{db: db}
}

func (sm *sensorManager) Add(sns *Sensor, dvc string) (string, error) {
	var json, err = json.Marshal(sns)
	if err != nil {
		return "", err
	}
	var key = fmt.Sprint(uuid.NewV4())
	err = sm.db.cli.Set(fmt.Sprintf("sensors:%v:device:%v",
		key, dvc), json, 0).Err()
	return key, err
}

func (sm *sensorManager) Read(key, dvc string) (sns Sensor, err error) {
	val, err := sm.db.cli.Get(fmt.Sprintf("sensors:%v:device:%v",
		key, dvc)).Result()
	if err != nil {
		return Sensor{}, err
	}

	err = json.Unmarshal([]byte(val), &sns)
	return sns, err
}

func (sm *sensorManager) Update(key string, dvc string, sns *Sensor) error {
	var json, err = json.Marshal(sns)
	if err != nil {
		return err
	}
	return sm.db.cli.Set(fmt.Sprintf("sensors:%v:device:%v",
		key, dvc), json, 0).Err()
}

func (sm *sensorManager) Delete(key, dvc string) error {
	return sm.db.cli.Del(fmt.Sprintf("sensors:%v:device:%v", key, dvc)).Err()
}

func (sm *sensorManager) ListByDevice(dvc string) ([]string, []Sensor, error) {
	var keys, _, err = sm.db.cli.Scan(0, "sensors:*:device:"+dvc, 0).Result()
	if err != nil {
		return []string{}, []Sensor{}, err
	}
	arr, err := sm.db.cli.MGet(keys...).Result()
	if err != nil {
		return []string{}, []Sensor{}, err
	}

	var sensors []Sensor
	var ids []string

	for i, elem := range arr {
		var sns Sensor
		err = json.Unmarshal([]byte(elem.(string)), &sns)
		if err != nil {
			return ids, sensors, err
		}

		ids = append(ids, strings.Split(keys[i], ":")[1])
		sensors = append(sensors, sns)
	}
	return ids, sensors, err
}