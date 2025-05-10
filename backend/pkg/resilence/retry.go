package resilence

import (
	"fmt"
	"math"
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(baseDelay) * time.Millisecond * time.Duration(math.Pow(2, float64(i))))
	}
	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, err)
}
