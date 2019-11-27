package easylogger

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/easylogger"
	"github.com/gomodule/redigo/redis"
)

var (
	configsPath                 string
	isPackageSetup              bool
	redisService, namingService *easylogger.Address
)

// LogLevel represents the how serious a log is
type LogLevel int

const (
	// DEBUG represents a debug log level
	DEBUG LogLevel = iota
	// INFO represents a informational log level
	INFO
	// WARNING represents a warning log level
	WARNING
	// ERROR represents an error log level
	ERROR
	// FATAL represents a fatal log level
	FATAL
)

// InitLogger initializes logger with configuration file
func InitLogger(loggerConfigsPath string) error {
	if isPackageSetup {
		return errors.New("InitLogger already called. It is not idempotent")
	}

	// Adds two more CPUs that this middleware requires to run in max performance
	runtime.GOMAXPROCS(runtime.NumCPU() + 2)
	configsPath = loggerConfigsPath
	configs := &easylogger.Configuration{}
	err := configuration.LoadConfiguration(loggerConfigsPath, &configs)
	if err != nil {
		return err
	}

	redisService = configs.RedisService
	namingService = configs.NamingService

	conn, err := redis.Dial("tcp", redisService.FullAddress())
	if err != nil {
		return err
	}

	res, err := redis.DoWithTimeout(conn, time.Second*2, "PING", "Hello, Redis!")
	if err != nil {
		return err
	}

	fmt.Printf("Redis server on, PING response: %+v\n", string(res.([]uint8)))
	isPackageSetup = true
	return nil

}

func log(message, destination, serviceID string, level LogLevel) {
	fmt.Printf("%s to %s as %d from %s", message, destination, level, serviceID)
}

// Log is the main function for logging
func Log(message, destination, serviceID string, level LogLevel) error {
	if !isPackageSetup {
		return errors.New("InitLogger hasn't been called yet or an error occurred on last call. Make sure EasyLogger package is correctly setup by calling it")
	}
	go log(message, destination, serviceID, level)
	return nil
}
