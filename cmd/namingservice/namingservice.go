package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration"
	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration/namingservice"
	nsserver "github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/namingservice"
)

func main() {
	numArgs := len(os.Args) - 1
	if numArgs != 3 {
		fmt.Printf("Wrong number of arguments. Expected 3 (ConfigFilePath, Host, Port), got %d", numArgs)
		os.Exit(1)
	}

	configs := &nsconfigs.Configuration{}
	err := configuration.LoadConfiguration(os.Args[1], configs)
	if err != nil {
		panic(err)
	}

	port64, err := strconv.ParseInt(os.Args[3], 10, 32)
	if err != nil {
		panic(err)
	}
	port := int(port64)

	for _, config := range configs.Loggers {
		fmt.Printf("%+v\n", config)
	}

	server := nsserver.InitNamingService(os.Args[2], port, configs)
	server.Start(2)

	time.Sleep(time.Hour)
}
