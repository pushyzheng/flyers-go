package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	FilenameEmptyError = errors.New("filename cannot be empty")
)

func ExistsFile(f string) bool {
	if len(f) == 0 {
		return false
	}
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}

func RemoveFile(f string) error {
	if len(f) == 0 {
		return FilenameEmptyError
	}
	if !ExistsFile(f) {
		return fileNotExists(f)
	}
	return os.Remove(f)
}

func ReadText(f string) (string, error) {
	if len(f) == 0 {
		return "", FilenameEmptyError
	}
	if !ExistsFile(f) {
		return "", fileNotExists(f)
	}
	b, err := os.ReadFile(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WriteText(f string, data string) error {
	if len(f) == 0 {
		return FilenameEmptyError
	}
	return os.WriteFile(f, []byte(data), 0666)
}

func ReadJsonText(f string) (map[string]any, error) {
	if len(f) == 0 {
		return nil, FilenameEmptyError
	}
	data := make(map[string]any)
	if !ExistsFile(f) {
		return data, nil
	}
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func fileNotExists(f string) error {
	return fmt.Errorf("file don't exists: %s", f)
}
