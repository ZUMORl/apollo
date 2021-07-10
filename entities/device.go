package entities

import (
	"encoding/json"
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

func (dev Device) Create(cli *redis.Client) {
	json, err := json.Marshal(dev)
	if err != nil {
		fmt.Println(err)
	}

	err = cli.Set(fmt.Sprintf("devices:%v", dev.id), json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func (dev Device) Update(cli *redis.Client) {

}

func (dev Device) Delete(cli *redis.Client) {

}

func (dev Device) Get(cli *redis.Client) {

}
