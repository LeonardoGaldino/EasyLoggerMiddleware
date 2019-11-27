package easylogger

import "fmt"

// Address is a struct containing a address information, e.g., host and port.
type Address struct {
	Host string
	Port int
}

// FullAddress returns full address as string (host:port)
func (addr *Address) FullAddress() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

// Configuration is a struct containing addresses information
type Configuration struct {
	RedisService, NamingService *Address
}
