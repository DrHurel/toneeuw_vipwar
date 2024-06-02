package cache

import (
	"auth"
	"encoding/json"
	"fmt"
	"os"
	"vipwar"
)

type Cache struct {
	VIP  vipwar.VIP `json:"vip"`
	User auth.User  `json:"user"`
}

func New(path string) (*Cache, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	res := new(Cache)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(res)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return res, nil
}

func (this *Cache) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(this)
	if err != nil {
		fmt.Println("Error encoding JSON:",
			err)
		return err
	}

	return nil
}
