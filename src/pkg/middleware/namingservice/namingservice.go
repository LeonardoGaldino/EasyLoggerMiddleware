package namingservice

import (
	"errors"
	"fmt"
	"net"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	nsMarshaller "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/marshaller/namingservice"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/network"
)

// NamingService is a struct for reaching the NamingService with easy API
type NamingService struct {
	address *configuration.Address
}

func (ns *NamingService) requestor(req *nsMarshaller.RequestMessage) (*nsMarshaller.ResponseMessage, error) {
	conn, err := net.Dial("tcp", ns.address.FullAddress())
	if err != nil {
		return nil, err
	}

	serializedReq := nsMarshaller.MarshallRequest(req)
	err = network.WriteMessage(&conn, serializedReq)
	if err != nil {
		return nil, err
	}

	resBytes, err := network.ReadMessage(&conn)
	if err != nil {
		return nil, err
	}

	response := nsMarshaller.UnmarshallResponse(resBytes)
	return response, nil
}

// Query queries a serviceName for its address
func (ns *NamingService) Query(serviceName string) (string, error) {
	req := &nsMarshaller.RequestMessage{
		Op:   nsMarshaller.QUERY,
		Data: serviceName,
	}
	response, err := ns.requestor(req)
	if err != nil {
		return "", err
	}
	if response.Res == nsMarshaller.OK {
		return response.Data, nil
	} else if response.Res == nsMarshaller.ERROR {
		return "", errors.New("Error on NamingService server")
	}
	return "", fmt.Errorf("Service %s not found on NamingService", serviceName)
}

// Register registers a service on the namingService
func (ns *NamingService) Register(serviceName, host string, port int) error {
	req := &nsMarshaller.RequestMessage{
		Op:   nsMarshaller.REGISTER,
		Data: fmt.Sprintf("%s:%s:%d", serviceName, host, port),
	}
	response, err := ns.requestor(req)
	if err != nil {
		return err
	}
	if response.Res == nsMarshaller.OK {
		return nil
	}
	return errors.New("Error on NamingService server")
}

// Unregister unregisters a service on the namingService
func (ns *NamingService) Unregister(serviceName string) error {
	req := &nsMarshaller.RequestMessage{
		Op:   nsMarshaller.UNREGISTER,
		Data: serviceName,
	}
	response, err := ns.requestor(req)
	if err != nil {
		return err
	}
	if response.Res == nsMarshaller.OK {
		return nil
	}
	return errors.New("Error on NamingService server")
}

// InitNamingService correctly initializes the struct NamingService from host and port
func InitNamingService(host string, port int) *NamingService {
	return &NamingService{
		address: &configuration.Address{
			Host: host,
			Port: port,
		},
	}
}

// InitNamingServiceFromAddr correctly initializes the struct NamingService from its adress
func InitNamingServiceFromAddr(addr *configuration.Address) *NamingService {
	return &NamingService{
		address: addr,
	}
}
