package easylogger

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/easylogger"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/network"
	nsMarshaller "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/network/marshaller/namingservice"
	"github.com/gomodule/redigo/redis"
)

var (
	configsPath                 string
	isPackageSetup              bool
	redisService, namingService *configuration.Address
	connPool                    *redis.Pool
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
	conn, err := net.Dial("tcp", namingService.FullAddress())
	if err != nil {
		return "", err
	}

	req := &nsMarshaller.RequestMessage{
		Op:   nsMarshaller.QUERY,
		Data: serviceName,
	}
	serializedReq := nsMarshaller.MarshallRequest(req)
	err = network.WriteMessage(&conn, serializedReq)
	if err != nil {
		return "", err
	}

	resBytes, err := network.ReadMessage(&conn)
	if err != nil {
		return "", err
	}

	response := nsMarshaller.UnmarshallResponse(resBytes)
	if response.Res == nsMarshaller.OK {
		return response.Data, nil
	} else if response.Res == nsMarshaller.ERROR {
		return "", errors.New("Error on NamingService server")
	}
	return "", fmt.Errorf("Service %s not found on NamingService", serviceName)
}

func logger(conn *redis.Conn) {
	pubsub := &redis.PubSubConn{Conn: *conn}
	pubsub.Subscribe("easylogger:logs")
	for {
		switch msg := pubsub.Receive().(type) {
		case redis.Message:
			data := string(msg.Data)
			fmt.Printf("Received: %s\n", data)
			addr, err := getServiceAddress(strings.Split(data, ":")[0])
			if err == nil {
				fmt.Printf("%s\n", addr)
			} else {
				fmt.Printf("Error: %+v\n", err)
			}
		}
	}
}

func initConnPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisService.FullAddress())
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
	runtime.GOMAXPROCS(runtime.NumCPU() + 2)

	configsPath = loggerConfigsPath
	configs := &easylogger.Configuration{}
	err := configuration.LoadConfiguration(loggerConfigsPath, &configs)
	if err != nil {
		return err
	}

	redisService = configs.RedisService
	namingService = configs.NamingService
	connPool = initConnPool()

	conn := connPool.Get()
	res, err := redis.DoWithTimeout(conn, time.Second*2, "PING", "Hello, Redis!")
	if err != nil {
		return err
	}

	fmt.Printf("Redis server on, PING response: %+v\n", string(res.([]uint8)))
	go logger(&conn)
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
