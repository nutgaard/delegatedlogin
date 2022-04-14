package internal

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	t.Log("Starting test")
	t.Log(3 * time.Second)
	execCounter := 0
	retryErr := Retry(3, 100*time.Millisecond, func() error {
		t.Log("Running retry")
		execCounter++
		return errors.New("something wrong")
	})

	t.Log(retryErr.Error())
	if execCounter != 3 {
		t.Fatalf("retry(3, 1, func) was executed %d times, expected %d", execCounter, 3)
	}
}
