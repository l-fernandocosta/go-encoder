package utils

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

func IsJson(s string) error {
	var js struct{}

	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return err
	}

	return nil
}

func GenerateUUIDString() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return id.String(), err
}
