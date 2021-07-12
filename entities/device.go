package entities

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

type Device struct {
	id    int
	Name  string
	Model string
}

func NewDevice(id int, param ...string) *Device {
	var dev = new(Device)
	dev.id = id
	if len(param) > 0 {
		dev.Name = param[0]
		dev.Model = param[1]
	}
	return dev
}

func (dev Device) Create(cli *redis.Client) (err error) {
	json, err := json.Marshal(dev)
	if err != nil {
		return err
	}

	err = cli.Set(fmt.Sprintf("devices:%v", dev.id), json, 0).Err()
	return err
}

func (dev Device) Update(cli *redis.Client) (err error) {
	return dev.Create(cli)
}

func (dev Device) Delete(cli *redis.Client) (err error) {
	return nil
}

func (dev *Device) Read(cli *redis.Client) (err error) {
	val, err := cli.Get(fmt.Sprintf("devices:%v", dev.id)).Result()
	switch err {
	case redis.Nil:
		return errors.New("such device not exist yet")
	case nil:
		break
	default:
		return err
	}

	err = json.Unmarshal([]byte(val), dev)
	return err
}
