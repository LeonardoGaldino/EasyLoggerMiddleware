package main

import (
	"fmt"
	"os"
	"time"

	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration/namingservice"
	nsserver "github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/namingservice"
)

func main() {
	numArgs := len(os.Args) - 1
	if numArgs != 1 {
		fmt.Printf("Wrong number of arguments. Expected 1 (config file path), got %d", numArgs)
		os.Exit(1)
	}

	configs := nsconfigs.LoadConfiguration(os.Args[1])

	for _, config := range configs.Loggers {
		fmt.Printf("%+v\n", config)
	}

	server := nsserver.InitNamingService("localhost", 8080, configs)
	server.Start(2)

	time.Sleep(time.Hour)
}
