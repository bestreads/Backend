package services

import "github.com/bestreads/Backend/internal/database"

func SaveFile(data []byte) (int, error) {
	return database.FileStoreRaw(data)
}
