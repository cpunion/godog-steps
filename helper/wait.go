package helper

import (
	"context"
	"fmt"
	"time"
)

func WaitFor(ctx context.Context, fn func() bool) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout")
		default:
			if fn() {
				return nil
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
