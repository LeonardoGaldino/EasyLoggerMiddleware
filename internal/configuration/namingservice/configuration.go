package namingservice

import (
	"encoding/json"
	"os"
)

// Logger is a struct containing a logger information, e.g., name and address.
type Logger struct {
	Name string
	Host string
	Port int
}

// Configuration is a struct containing loggers information
type Configuration struct {
	Loggers []*Logger
}

// LoadConfiguration loads Configuration struct from a configuration file
func LoadConfiguration(path string) *Configuration {
	file, _ := os.Open(path)
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	return &configuration
}
