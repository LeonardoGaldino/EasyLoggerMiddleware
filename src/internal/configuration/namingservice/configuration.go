package namingservice

import (
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
)

// Logger is a struct containing a logger information, e.g., name and address.
type Logger struct {
	Name    string
	Address *configuration.Address
}

// Configuration is a struct containing loggers information
type Configuration struct {
	SelfAddress *configuration.Address
	Loggers     []*Logger
}
