package namingservice

import (
	"fmt"
	"io"
	"net"

	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/namingservice"
)

// NamingService represents a namingService address and its data
type NamingService struct {
	Host string
	Port int
	data map[string]string
}

func (service *NamingService) handle(conn *net.Conn) {
	bufferSize := 512
	buffer := make([]byte, bufferSize, bufferSize)

	messageSize, err := (*conn).Read(buffer)
	if err != nil && err != io.EOF {
		panic(err)
	}
	message := string(buffer[:messageSize])
	fmt.Println(message)

	(*conn).Write([]byte(service.data[message]))
	(*conn).Close()
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
		data[config.Name] = fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	return &NamingService{
		Host: host,
		Port: port,
		data: data,
	}
}
