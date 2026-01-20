package services

import "github.com/bestreads/Backend/internal/database"

func SaveFile(data []byte) (uint64, error) {
	return database.FileStoreRaw(data)
}
