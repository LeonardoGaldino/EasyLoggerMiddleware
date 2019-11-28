package namingservice

import (
	"fmt"
	"net"

	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/namingservice"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/network"
)

// NamingService represents a namingService address and its data
type NamingService struct {
	Host string
	Port int
	data map[string]string
}

func (service *NamingService) handle(conn *net.Conn) {
	var res string
	buffer, err := network.ReadMessage(conn)
	if err != nil {
		fmt.Printf("Error on receiving message: %+v\n", err)
		res = "FAIL/"
		network.WriteMessage(conn, []byte(res))
	} else {
		message := string(buffer)
		fmt.Printf("%s\n", message)
		res = service.data[message]
		if len(res) > 0 {
			res = fmt.Sprintf("OK/%s", res)
		} else {
			res = "FAIL/"
		}
	}
	network.WriteMessage(conn, []byte(res))
}

// Start starts NamingService TCP server
func (service *NamingService) Start(maxConcurrency int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", service.Host, service.Port))
	if err != nil {
		panic(err)
	}

	for i := 0; i < maxConcurrency; i++ {
		go func(id int) {
			for {
				conn, err := listener.Accept()
				if err != nil {
					panic(err)
				}
				service.handle(&conn)
				conn.Close()
			}
		}(i)
	}

}

// InitNamingService initializes NamingService struct from host,port and configs.
func InitNamingService(host string, port int, configs *nsconfigs.Configuration) *NamingService {
	data := make(map[string]string)
	for _, config := range configs.Loggers {
		data[config.Name] = config.Address.FullAddress()
	}

	return &NamingService{
		Host: host,
		Port: port,
		data: data,
	}
}
