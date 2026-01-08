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

var ImageTypeMap = map[int]ImageType{
	0: PostImage,
	1: ProfileImage,
}

func FileStoreRaw(data []byte, itype ImageType) (int, error) {
	val, err := fnv.New128a().Write(data)
	if err != nil {
		return -1, err
	}

	err = os.WriteFile(prefix(strconv.Itoa(val), itype), data, 0640)
	if err != nil {
		return -1, err
	}

	return val, nil

}

func FileStoreB64(data string, itype ImageType) (int, error) {
	// schneller pfad, kein fs-aufruf
	if data == "" {
		return 0, nil
	}

	bytes := []byte(data)
	return FileStoreRaw(bytes, itype)
}

func FileRetrieve(hash string, itype ImageType) ([]byte, error) {
	return os.ReadFile(prefix(hash, itype))
}

func FileRetrieveB64(hash string, itype ImageType) (string, error) {
	// auch schneller pfad
	if hash == "0" {
		return "", nil
	}
	d, err := FileRetrieve(hash, itype)
	if err != nil {
		return "", err
	}

	return string(d), nil

}

func prefix(name string, itype ImageType) string {
	// eigentlich müsssen wir hier noch den pfad sanitizen
	return fmt.Sprintf("./store/%d/%s", itype, name)
}
