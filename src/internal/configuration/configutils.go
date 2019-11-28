package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

// Address is a struct containing a address information, e.g., host and port.
type Address struct {
	Host string
	Port int
}

// FullAddress returns full address as string (host:port)
func (addr *Address) FullAddress() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

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
