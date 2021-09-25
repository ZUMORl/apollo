package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type (
	Sensor struct {
		Type  string `json:"type"`
		Model string `json:"model"`
	}

	Sensors interface {
		Add(*Sensor, string) (string, error)
		Read(string) (Sensor, error)
		Update(string, *Sensor) error
		Delete(string) error
		AddValue(string, string) error
		RemoveValue(string, int, int) error
		GetValues(string, int, int) ([]string, error)
		ListByDevice(string) (map[string]Sensor, error)
	}

	sensorManager struct {
		db DataBase
	}
)

func getFullKey(id string, sm *sensorManager) (string, error) {
	var keys, _, err = sm.db.cli.Scan(0, "sensors:"+id+"*", 0).Result()
	if err != nil {
		return "", err
	}
	if len(keys) == 1 {
		return keys[0], nil
	}
	return "", errors.New("id is not unique")
}

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

func (sm *sensorManager) Read(id string) (sns Sensor, err error) {
	key, err := getFullKey(id, sm)
	if err != nil {
		return Sensor{}, err
	}
	val, err := sm.db.cli.Get(key).Result()
	if err != nil {
		return Sensor{}, err
	}

	err = json.Unmarshal([]byte(val), &sns)
	return sns, err
}

func (sm *sensorManager) Update(id string, sns *Sensor) error {
	var json, err = json.Marshal(sns)
	if err != nil {
		return err
	}
	key, err := getFullKey(id, sm)
	if err != nil {
		return err
	}
	return sm.db.cli.Set(key, json, 0).Err()
}

func (sm *sensorManager) Delete(id string) error {
	key, err := getFullKey(id, sm)
	if err != nil {
		return err
	}
	return sm.db.cli.Del(key).Err()
}

func (sm *sensorManager) ListByDevice(dvc string) (map[string]Sensor, error) {
	var keys, _, err = sm.db.cli.Scan(0, "sensors:*:device:"+dvc, 0).Result()
	if err != nil || len(keys) == 0 {
		return map[string]Sensor{}, err
	}
	arr, err := sm.db.cli.MGet(keys...).Result()
	if err != nil {
		return map[string]Sensor{}, err
	}

	var ret = map[string]Sensor{}

	for i, elem := range arr {
		var sns Sensor
		err = json.Unmarshal([]byte(elem.(string)), &sns)
		if err != nil {
			return ret, err
		}

		ret[strings.Split(keys[i], ":")[1]] = sns
	}
	return ret, err
}

func (sm *sensorManager) AddValue(id string, value string) error {
	var key, err = getFullKey(id, sm)
	if err != nil {
		return err
	}

	return sm.db.cli.LPush("values:"+key, value).Err()
}

func (sm *sensorManager) RemoveValue(id string, start int, end int) error {
	var key, err = getFullKey(id, sm)
	if err != nil {
		return err
	}
	return sm.db.cli.LTrim(key, int64(start), int64(end)).Err()
}

func (sm *sensorManager) GetValues(id string, start int, end int) ([]string, error) {
	var key, err = getFullKey(id, sm)
	if err != nil {
		return []string{}, err
	}

	return sm.db.cli.LRange(key, int64(start), int64(end)).Result()
}
