package dispatcher

import (
	"fmt"
	"strings"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/utils"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	core "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/pkg/middleware/easylogger"
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
	pubsub.Subscribe(core.RedisChannel)
	for {
		switch msg := pubsub.Receive().(type) {
		case redis.Message:
			data := string(msg.Data)
			fmt.Printf("Received: %s\n", data)
			fields := strings.Split(data, ":")
			addr := utils.KeepRetryingAfter(func() (interface{}, error) {
				addr, err := namingService.Query(fields[0])
				if err != nil && strings.Contains(err.Error(), "not found") {
					return addr, nil
				}
				return addr, err
			}, time.Second)
			fmt.Printf("%s\n", addr)
			/*
			 * TODO: demultiplex using the destination (fields[0]) and the address (addr)
			 * and send the log content fields[1:] to log service
			 */
		}
	}
}
