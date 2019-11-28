package easylogger

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/easylogger"
	nsAPI "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/pkg/middleware/namingservice"
	"github.com/gomodule/redigo/redis"
)

var (
	configsPath                         string
	isPackageSetup                      bool
	redisServiceAddr, namingServiceAddr *configuration.Address
	namingService                       *nsAPI.NamingService
	connPool                            *redis.Pool
)

// LogLevel represents the how serious a log is
type LogLevel int

func (l LogLevel) String() string {
	return []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}[l]
}

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

func getServiceAddress(serviceName string) (string, error) {
	return namingService.Query(serviceName)
}

func initConnPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisServiceAddr.FullAddress())
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxActive: 11,
		Wait:      true,
	}
}

// InitLogger initializes logger with configuration file
func InitLogger(loggerConfigsPath string) error {
	if isPackageSetup {
		return errors.New("InitLogger already called. It is not idempotent")
	}

	// Adds two more CPUs that this middleware requires to run in max performance
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)

	configsPath = loggerConfigsPath
	configs := &easylogger.Configuration{}
	err := configuration.LoadConfiguration(loggerConfigsPath, &configs)
	if err != nil {
		return err
	}

	redisServiceAddr = configs.RedisService
	namingServiceAddr = configs.NamingService
	namingService = nsAPI.InitNamingServiceFromAddr(namingServiceAddr)
	connPool = initConnPool()

	conn := connPool.Get()
	res, err := redis.DoWithTimeout(conn, time.Second*2, "PING", "Hello, Redis!")
	if err != nil {
		return err
	}

	fmt.Printf("Redis server on, PING response: %+v\n", string(res.([]uint8)))
	isPackageSetup = true
	return nil
}

func keepRetryingAfter(f func() (interface{}, error), after time.Duration) interface{} {
	v, err := f()
	for {
		if err == nil {
			break
		}
		time.Sleep(after)
		v, err = f()
	}
	return v
}

func log(message, destination, serviceID string, level LogLevel) {
	conn := connPool.Get()
	defer conn.Close()

	now := time.Now().Unix()
	serialized := fmt.Sprintf("%s:%d:%s:%s:%s", destination, now, level.String(), serviceID, message)

	keepRetryingAfter(func() (interface{}, error) {
		return conn.Do("PUBLISH", "easylogger:logs", serialized)
	}, time.Second*3)
}

// Log is the main function for logging
func Log(message, destination, serviceID string, level LogLevel) error {
	if !isPackageSetup {
		return errors.New("InitLogger hasn't been called yet or an error occurred on last call. Make sure EasyLogger package is correctly setup by calling it")
	}
	go log(message, destination, serviceID, level)
	return nil
}
