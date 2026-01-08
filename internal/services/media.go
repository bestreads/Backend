package services

import "github.com/bestreads/Backend/internal/database"

func SaveFile(data []byte, itype database.ImageType) (int, error) {
	return database.FileStoreRaw(data, itype)
}
