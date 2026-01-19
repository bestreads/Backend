package database

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"strconv"
)

// speichert die gegebene daten im filestore. gibt den hash (schlüssel) wieder zurück
func FileStoreRaw(data []byte) (uint64, error) {
	h := fnv.New64a()
	h.Write(data) // Write auf hash.Hash gibt nie einen Fehler zurück
	hash := h.Sum64()

	err := os.WriteFile(prefix(strconv.FormatUint(hash, 10)), data, 0640)
	if err != nil {
		return 0, err
	}

	return hash, nil
}

// Deprecated: mach mal nicht
func FileStoreB64(data string) (uint64, error) {
	// schneller pfad, kein fs-aufruf
	if data == "" {
		return 0, nil
	}

	bytes := []byte(data)
	return FileStoreRaw(bytes)
}

// sucht nach dem hash (schlüssel) und gibt die datei wieder zurück.
func FileRetrieve(hash string) ([]byte, error) {
	return os.ReadFile(prefix(hash))
}

// Deprecated: mach mal nicht
func FileRetrieveB64(hash string) (string, error) {
	// auch schneller pfad
	if hash == "0" {
		return "", nil
	}
	d, err := FileRetrieve(hash)
	if err != nil {
		return "", err
	}

	return string(d), nil

}

func prefix(name string) string {
	// eigentlich müsssen wir hier noch den pfad sanitizen
	return fmt.Sprintf("./store/%s", name)
}

// cached den link im dateisystem. gibt den hash (schlüssel) wieder zurück
func CacheMedia(url string) (uint64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	hash, err := FileStoreRaw(body)
	if err != nil {
		return 0, err
	}

	return hash, nil
}
