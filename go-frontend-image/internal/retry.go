package internal

import "time"

func Retry(attempts int, sleepDuration time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleepDuration)
			return Retry(attempts, sleepDuration, fn)
		}
		return err
	}
	return nil
}
