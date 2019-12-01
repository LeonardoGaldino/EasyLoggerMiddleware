package utils

import (
	"math"
	"sort"
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

// SortableDurationSlice is a struct for sorting duration slice
type SortableDurationSlice struct {
	Data []time.Duration
}

// Len returns the size of underlying duration slice
func (s SortableDurationSlice) Len() int {
	return len(s.Data)
}

// Less returns a boolean indicating if element at i is less then element at j
func (s SortableDurationSlice) Less(i, j int) bool {
	return (s.Data[i]) < (s.Data[j])
}

// Swap swaps two element in the underlying slice
func (s SortableDurationSlice) Swap(i, j int) {
	temp := s.Data[i]
	s.Data[i] = s.Data[j]
	s.Data[j] = temp
}

// ComputeMetrics computes avg and standard deviation for a set of durations (microseconds)
func ComputeMetrics(delays []time.Duration) (float64, float64, int, time.Duration) {
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

	sortable := SortableDurationSlice{Data: delays}
	sort.Sort(sortable)
	var p90index = int(math.Floor(0.90 * float64(sortable.Len())))
	return avgNanoSecs / 1000, sd / 1000, zeroValues, sortable.Data[p90index]
}
