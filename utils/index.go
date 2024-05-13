package utils

import (
	"os"
	"time"

	"github.com/oklog/ulid/v2"
)


func GenerateUUIDString() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulid.DefaultEntropy()).String()
}

func LoadFile(filePath string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}