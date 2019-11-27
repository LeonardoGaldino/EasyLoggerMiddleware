package easylogger

import (
	"fmt"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration/easylogger"
	"github.com/gomodule/redigo/redis"
)

var (
	configsPath                 string
	isPackageSetup              bool
	redisService, namingService *easylogger.Address
)

// InitLogger initializes logger with configuration file
func InitLogger(loggerConfigsPath string) error {
	configsPath = loggerConfigsPath
	configs := &easylogger.Configuration{}
	err := configuration.LoadConfiguration(loggerConfigsPath, &configs)
	if err != nil {
		return err
	}

	redisService = configs.RedisService
	namingService = configs.NamingService
	isPackageSetup = true

	conn, err := redis.Dial("tcp", redisService.FullAddress())
	if err != nil {
		return err
	}

	res, err := redis.DoWithTimeout(conn, time.Second*2, "PING", "Hello, Redis!")
	if err != nil {
		return err
	}

	fmt.Printf("Redis server on, PING response: %+v\n", string(res.([]uint8)))
	return nil
}
