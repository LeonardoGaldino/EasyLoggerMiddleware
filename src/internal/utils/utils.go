package utils

import (
	"math"
	"time"
)

// KeepRetryingAfter keeps retring to call some function over and over after 'after' duration.
func KeepRetryingAfter(f func() (interface{}, error), after time.Duration) interface{} {
	v, err := f()
	for {
		if err == nil {
			break
		}
		time.Sleep(after)
		v, err = f()
	}
	return v
}

// ComputeMetrics computes avg and standard deviation for a set of durations (microseconds)
func ComputeMetrics(delays []time.Duration) (float64, float64, int) {
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
	return avgNanoSecs / 1000, sd / 1000, zeroValues
}
