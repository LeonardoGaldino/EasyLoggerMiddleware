package main

import (
	"fmt"
	"os"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration"
	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration/namingservice"
	nsserver "github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/namingservice"
)

func main() {
	numArgs := len(os.Args) - 1
	if numArgs != 1 {
		fmt.Printf("Wrong number of arguments. Expected 1 (config file path), got %d", numArgs)
		os.Exit(1)
	}

	configs := &nsconfigs.Configuration{}
	err := configuration.LoadConfiguration(os.Args[1], configs)
	if err != nil {
		panic(err)
	}

	for _, config := range configs.Loggers {
		fmt.Printf("%+v\n", config)
	}

	server := nsserver.InitNamingService("localhost", 8080, configs)
	server.Start(2)

	time.Sleep(time.Hour)
}
