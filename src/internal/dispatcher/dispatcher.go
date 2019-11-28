package dispatcher

import (
	"fmt"
	"strings"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	nsAPI "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/pkg/middleware/namingservice"
	"github.com/gomodule/redigo/redis"
)

// StartDispatching subscribe on Redis channel and waits for logs to get them and log it to destination
func StartDispatching(redisAddr *configuration.Address, namingServiceAddr *configuration.Address) error {
	conn, err := redis.Dial("tcp", redisAddr.FullAddress())
	if err != nil {
		return err
	}
	namingService := nsAPI.InitNamingServiceFromAddr(namingServiceAddr)

	pubsub := &redis.PubSubConn{Conn: conn}
	pubsub.Subscribe("easylogger:logs")
	for {
		switch msg := pubsub.Receive().(type) {
		case redis.Message:
			data := string(msg.Data)
			fmt.Printf("Received: %s\n", data)
			fields := strings.Split(data, ":")
			addr, err := namingService.Query(fields[0])
			if err == nil {
				fmt.Printf("%s\n", addr)
			} else {
				fmt.Printf("Error: %+v\n", err)
			}
		}
	}
}
