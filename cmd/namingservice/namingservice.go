package main

import (
	"fmt"
	"os"

	nsconfigs "github.com/LeonardoGaldino/EasyLoggerMiddleware/internal/configuration/namingservice"
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
}
