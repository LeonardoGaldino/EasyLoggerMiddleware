package namingservice

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
