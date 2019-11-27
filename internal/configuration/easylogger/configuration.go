package easylogger

// Address is a struct containing a address information, e.g., host and port.
type Address struct {
	Host string
	Port int
}

// Configuration is a struct containing addresses information
type Configuration struct {
	RedisService, NamingService *Address
}
