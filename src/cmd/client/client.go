package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/pkg/middleware/easylogger"
)

func main() {
	args := os.Args

	easylogger.InitLogger("../../../easylogger_configs.json")
	reps, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Print("enter the number of times you want to log as argument")
		panic(err)
	}
	var totalElapsed []time.Duration

	for i := 0; i < reps; i++ {
		start := time.Now()
		easylogger.Log("xalxala", "elasticsearch", "my_client", 0)
		elapsed := time.Since(start)
		totalElapsed = append(totalElapsed, elapsed)
	}

	avg, sd := computeMetrics(totalElapsed)
	println(avg, "ns", sd, "ns")
}

func computeMetrics(delays []time.Duration) (float64, float64) {
	var totalNanoSecs int64
	len := float64(len(delays))
	zeroValues := 0
	for _, delay := range delays {
		if delay == 0 {
			zeroValues++
		}
		totalNanoSecs += delay.Nanoseconds()
	}
	avgNanoSecs := float64(totalNanoSecs) / (len - float64(zeroValues))

	var sd float64
	for _, delay := range delays {
		diff := (float64(delay.Nanoseconds()) - avgNanoSecs)
		sd += diff * diff
	}
	sd = sd / (len - 1 - float64(zeroValues))
	sd = math.Sqrt(sd)
	return avgNanoSecs / 1000, sd / 1000
}
