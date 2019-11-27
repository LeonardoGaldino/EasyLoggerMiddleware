package easylogger

import (
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration/easylogger"
)

var (
	configsPath                 string
	isPackageSetup              bool
	redisService, namingService *easylogger.Address
)

// InitLogger initializes logger with configuration file
func InitLogger(loggerConfigsPath string) {
	configsPath = loggerConfigsPath
	configs := &easylogger.Configuration{}
	configuration.LoadConfiguration(loggerConfigsPath, &configs)
	redisService = configs.RedisService
	namingService = configs.NamingService
	isPackageSetup = true
}
