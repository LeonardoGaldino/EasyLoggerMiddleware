package persistence

import (
	"encoding/json"
	"io"
	"os"
)

// MarshallEntries marshalls entries into a string that can be written to file
func MarshallEntries(entries map[int]string) (string, error) {
	bytes, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// UnmarshallEntries unmarshalls entries from a file into a map
func UnmarshallEntries(file *os.File) (map[int]string, error) {
	decoder := json.NewDecoder(file)
	content := make(map[int]string)
	err := decoder.Decode(&content)
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	return content, nil
}
