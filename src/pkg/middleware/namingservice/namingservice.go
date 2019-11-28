package namingservice

import (
	"errors"
	"fmt"
	"net"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/network"
	nsMarshaller "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/network/marshaller/namingservice"
)

// NamingService is a struct for reaching the NamingService with easy API
type NamingService struct {
	address *configuration.Address
}

// Query queries a serviceName for its address
func (ns *NamingService) Query(serviceName string) (string, error) {
	conn, err := net.Dial("tcp", ns.address.FullAddress())
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

// Register registers a service on the namingService
func (ns *NamingService) Register(serviceName, host string, port int) error {
	conn, err := net.Dial("tcp", ns.address.FullAddress())
	if err != nil {
		return err
	}

	req := &nsMarshaller.RequestMessage{
		Op:   nsMarshaller.REGISTER,
		Data: fmt.Sprintf("%s:%s:%d", serviceName, host, port),
	}
	serializedReq := nsMarshaller.MarshallRequest(req)
	err = network.WriteMessage(&conn, serializedReq)
	if err != nil {
		return err
	}

	resBytes, err := network.ReadMessage(&conn)
	if err != nil {
		return err
	}

	response := nsMarshaller.UnmarshallResponse(resBytes)
	if response.Res == nsMarshaller.OK {
		return nil
	}
	return errors.New("Error on NamingService server")
}

// Unregister unregisters a service on the namingService
func (ns *NamingService) Unregister(serviceName string) error {
	conn, err := net.Dial("tcp", ns.address.FullAddress())
	if err != nil {
		return err
	}

	req := &nsMarshaller.RequestMessage{
		Op:   nsMarshaller.REGISTER,
		Data: serviceName,
	}
	serializedReq := nsMarshaller.MarshallRequest(req)
	err = network.WriteMessage(&conn, serializedReq)
	if err != nil {
		return err
	}

	resBytes, err := network.ReadMessage(&conn)
	if err != nil {
		return err
	}

	response := nsMarshaller.UnmarshallResponse(resBytes)
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
