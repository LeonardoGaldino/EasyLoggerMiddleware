package configuration

import (
	"encoding/json"
	"os"
)

// LoadConfiguration loads Configuration struct from a configuration file
func LoadConfiguration(path string, dest interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(dest)
	return err
}
