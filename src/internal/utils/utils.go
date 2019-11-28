package utils

import "time"

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
