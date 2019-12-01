package dispatcher

import (
	"fmt"
	"strings"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/persistence"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/utils"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	core "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/pkg/middleware/easylogger"
	nsAPI "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/pkg/middleware/namingservice"
	"github.com/gomodule/redigo/redis"
)

var (
	persistFile     = "dispatcher.persistence"
	persistor       = &persistence.Persistor{FileName: persistFile}
	dispatcherDemux = map[string]func(string, string){
		"ElasticSearch": elasticSearchDispatcher,
	}
)

// StartDispatching subscribe on Redis channel and waits for logs to get them and log it to destination
func StartDispatching(redisAddr *configuration.Address, namingServiceAddr *configuration.Address) error {
	conn, err := redis.Dial("tcp", redisAddr.FullAddress())
	if err != nil {
		return err
	}
	namingService := nsAPI.InitNamingServiceFromAddr(namingServiceAddr)
	persistor.GenericDispatchEntries(dispatcherDemux, namingService)
	pubsub := &redis.PubSubConn{Conn: conn}
	pubsub.Subscribe(core.RedisChannel)
	fmt.Printf("Starting to consume Redis channel: %s\n", core.RedisChannel)
	for {
		switch msg := pubsub.Receive().(type) {
		case redis.Message:
			data := string(msg.Data)
			fmt.Printf("Received: %s\n", data)
			id := persistor.AddEntry(data)
			fields := strings.Split(data, ":")

			addr := utils.KeepRetryingAfter(func() (interface{}, error) {
				addr, err := namingService.Query(fields[0])
				if err != nil && strings.Contains(err.Error(), "not found") {
					return addr, nil
				}
				return addr, err
			}, time.Second)
			fulladdr := fmt.Sprintf("http://%s", addr.(string))
			fmt.Printf("%s\n", fulladdr)

			dispatcher := dispatcherDemux[fields[0]]
			if dispatcher != nil {
				dispatcher(fulladdr, data)
			} else {
				fmt.Printf("No dispatcher for destination: %s\n", fields[0])
			}
			persistor.RemoveEntry(id)
		}
	}
}
