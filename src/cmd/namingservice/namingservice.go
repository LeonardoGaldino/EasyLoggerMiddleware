package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/namingservice"
	nsserver "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/namingservice"
)

var (
	expectedArgs       = []string{"ConfigFilePath"}
	numCommandLineArgs = len(expectedArgs)
)

func main() {
	numArgs := len(os.Args) - 1
	if numArgs != numCommandLineArgs {
		fmt.Printf("Wrong number of arguments. Expected 1 (%v) but got %d\n", expectedArgs, numArgs)
		os.Exit(1)
	}

	configs := &nsconfigs.Configuration{}
	err := configuration.LoadConfiguration(os.Args[1], configs)
	if err != nil {
		panic(err)
	}

	for _, config := range configs.Loggers {
		fmt.Printf("%s: %+v\n", config.Name, config.Address)
	}

	server := nsserver.InitNamingService(configs.SelfAddress.Host, configs.SelfAddress.Port, configs)
	server.Start(2)

	// Stop main goroutine while other goroutines handle incoming requests
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
