package main

import (
	"fmt"
	"os"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration"
	easyLoggerConfig "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/configuration/easylogger"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/dispatcher"
)

var (
	expectedArgs       = []string{"MiddlewareConfigsFilePath"}
	numCommandLineArgs = len(expectedArgs)
)

func main() {
	numArgs := len(os.Args) - 1
	if numArgs != numCommandLineArgs {
		fmt.Printf("Wrong number of arguments. Expected %d (%v) but got %d\n", numCommandLineArgs, expectedArgs, numArgs)
		os.Exit(1)
	}

	configs := &easyLoggerConfig.Configuration{}
	err := configuration.LoadConfiguration(os.Args[1], configs)
	if err != nil {
		panic(err)
	}

	err = dispatcher.StartDispatching(configs.RedisService, configs.NamingService)
	if err != nil {
		panic(err)
	}
}
