package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/utils"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/pkg/middleware/easylogger"
)

var (
	expectedArgs       = []string{"Configs", "RepetitionsNumber"}
	numCommandLineArgs = len(expectedArgs)
)

func main() {
	numArgs := len(os.Args) - 1
	if numArgs != numCommandLineArgs {
		fmt.Printf("Wrong number of arguments. Expected %d (%v) but got %d\n", numCommandLineArgs, expectedArgs, numArgs)
		os.Exit(1)
	}

	reps, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Print("enter the number of times you want to log as argument")
		panic(err)
	}

	easylogger.InitLogger(os.Args[1])

	var totalElapsed []time.Duration

	for i := 0; i < reps; i++ {
		start := time.Now()
		easylogger.Log("test", "ElasticSearch", "my_client", easylogger.FATAL)
		elapsed := time.Since(start)
		totalElapsed = append(totalElapsed, elapsed)
	}

	avg, sd, zeroes := utils.ComputeMetrics(totalElapsed)
	fmt.Printf("From a total of %d calls: [AVG: %.2f μs, SD: %.2f μs, %d zeroed-values]\n", reps, avg, sd, zeroes)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
