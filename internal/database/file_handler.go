package database

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var storeInitOnce sync.Once

// speichert die gegebene daten im filestore. gibt den hash (schlüssel) wieder zurück
func FileStoreRaw(data []byte) (uint64, error) {
	// Stelle sicher, dass der store-Ordner existiert
	if err := ensureStoreDir(); err != nil {
		return 0, fmt.Errorf("failed to create store directory: %w", err)
	}

	// ist das irgendwie nirgendswo dokumentiert?
	// WIE BENUTZT MAN DIE APIs?
	// https://gist.github.com/wjkoh/cd97f19cae5a9ac8a9fa61d4c6931b9e
	hash := fnv.New64a()
	_, err := hash.Write(data)
	if err != nil {
		return 0, err
	}

	// auf 128bit fnv hashes scheint die sum-funktion
	// einfach nen byte-array mit dem wert auszuspucken
	// daher jetzt einfach 64bit
	res := hash.Sum64()

	fmt.Printf("hash: %d", res)

	err = os.WriteFile(prefix(strconv.FormatUint(res, 10)), data, 0640)
	if err != nil {
		return 0, err
	}

	return res, nil
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

func getStorePath() string {
	// :) importzyklus von middlewares
	base := os.Getenv("DATA_PATH")
	if base == "" {
		base = "./store"
	}
	return base
}

func prefix(name string) string {
	// eigentlich müsssen wir hier noch den pfad sanitizen
	return fmt.Sprintf("%s/%s", getStorePath(), name)
}

// ensureStoreDir stellt sicher, dass der store-Ordner existiert
func ensureStoreDir() error {
	var err error
	storeInitOnce.Do(func() {
		err = os.MkdirAll(getStorePath(), 0755)
	})
	return err
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
