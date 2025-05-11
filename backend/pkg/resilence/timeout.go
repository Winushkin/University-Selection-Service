package resilence

import (
	"errors"
	"time"
)

// Timeout makes operation timed out
func Timeout(operation func() error, timeout int) error {
	done := make(chan error, 1)
	go func() {
		done <- operation()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		return errors.New("operation timed out")
	}
}
