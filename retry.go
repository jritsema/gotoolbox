package gotoolbox

import "time"

//Stop represents a stop error
type Stop struct {
	error
}

//Retry retries a function with exponential backoff
//return a Stop error to abort
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(Stop); ok {
			return s.error
		}
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}
