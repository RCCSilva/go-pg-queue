package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func Retry(
	ctx context.Context,
	attempts int,
	sleep time.Duration,
	max time.Duration,
	f func() error,
) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Printf("[retry] sleeping for %v", sleep)
			time.Sleep(sleep)
			sleep = min(sleep*2+(time.Duration(rand.Intn(1000))*time.Millisecond), max)
		}

		err = f()

		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("after %d attempts, failed with %v", attempts, err)
}

func min(x, y time.Duration) time.Duration {
	if x < y {
		return x
	}

	return y
}
