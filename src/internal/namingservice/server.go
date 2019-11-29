package namingservice

import (
	"fmt"
	"net"
	"strings"

	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/namingservice"
	nsMarshaller "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/marshaller/namingservice"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/network"
)

// NamingService represents a namingService address and its data
type NamingService struct {
	Host string
	Port int
	data map[string]string
}

// Register registers a service under the NamingService data
func (service *NamingService) Register(name, address, port string) {
	service.data[name] = fmt.Sprintf("%s:%s", address, port)
}

// Unregister removes a service of the NamingService data
func (service *NamingService) Unregister(name string) {
	delete(service.data, name)
}

// Query gets the address of the service running with a given name
func (service *NamingService) Query(name string) string {
	return service.data[name]
}

func (service *NamingService) demuxOperation(req *nsMarshaller.RequestMessage) *string {
	var res *string
	switch req.Op {
	case nsMarshaller.QUERY:
		temp := service.Query(req.Data)
		res = &temp
	case nsMarshaller.REGISTER:
		fields := strings.Split(req.Data, ":")
		service.Register(fields[0], fields[1], fields[2])
	case nsMarshaller.UNREGISTER:
		service.Unregister(req.Data)
	}
	return res
}

func (service *NamingService) handle(conn *net.Conn) {
	var data string
	var result nsMarshaller.Result
	buffer, err := network.ReadMessage(conn)
	if err != nil {
		fmt.Printf("Error on receiving message: %+v\n", err)
		result = nsMarshaller.ERROR
	} else {
		msg := nsMarshaller.UnmarshallRequest(buffer)
		res := service.demuxOperation(msg)
		if res == nil {
			result = nsMarshaller.OK
		} else if len(*res) > 0 {
			data = *res
			result = nsMarshaller.OK
		} else {
			result = nsMarshaller.NOTFOUND
		}
	}

	message := &nsMarshaller.ResponseMessage{
		Res:  result,
		Data: data,
	}
	serialized := nsMarshaller.MarshallResponse(message)
	network.WriteMessage(conn, serialized)
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
