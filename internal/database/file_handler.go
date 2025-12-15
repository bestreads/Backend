package database

import (
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
)

// sachen, um bilder zu verwalten

type ImageType int

const (
	PostImage ImageType = iota
	ProfileImage
)

func Store(data string, itype ImageType) (int, error) {
	bytes := []byte(data)
	val, err := fnv.New128a().Write(bytes)
	if err != nil {
		return -1, err
	}

	err = os.WriteFile(prefix(strconv.Itoa(val), itype), bytes, 0640)
	if err != nil {
		return -1, err
	}

	return val, nil
}

func prefix(name string, itype ImageType) string {

	return fmt.Sprintf("./store/%d/%s", itype, name)
}

func Retrieve(hash string, itype ImageType) (string, error) {
	d, err := os.ReadFile(prefix(hash, itype))
	if err != nil {
		return "", err
	}

	return string(d), nil

}
