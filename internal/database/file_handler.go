package database

import (
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
)

// sachen, um bilder zu verwalten

func Store(data string) (int, error) {
	bytes := []byte(data)
	val, err := fnv.New128a().Write(bytes)
	if err != nil {
		return -1, err
	}

	err = os.WriteFile(prefix(strconv.Itoa(val)), bytes, 0640)
	if err != nil {
		return -1, err
	}

	return val, nil
}

func prefix(name string) string {
	return fmt.Sprintf("./store/%s", name)
}

func Retrieve(hash string) (string, error) {
	d, err := os.ReadFile(prefix(hash))
	if err != nil {
		return "", err
	}

	return string(d), nil

}
