package easylogger

import (
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
)

// Configuration is a struct containing addresses information
type Configuration struct {
	RedisService, NamingService *configuration.Address
}
