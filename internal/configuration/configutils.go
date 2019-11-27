package configuration

import (
	"encoding/json"
	"os"
)

// LoadConfiguration loads Configuration struct from a configuration file
func LoadConfiguration(path string, dest interface{}) {
	file, _ := os.Open(path)
	defer file.Close()

	decoder := json.NewDecoder(file)

	err := decoder.Decode(dest)
	if err != nil {
		panic(err)
	}
}
