package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type RetryOptions struct {
	RetryCount    int
	WaitSeconds   int    // this is the nominal  wait seconds
	RetryStrategy string // [linear*, beb] // beb is binary exponential backoff
}

type RunnerFunc func(args ...interface{}) (any, error)

// RunWithRetry retries the provided function with the given arguments until it succeeds or runs out of retries.
// It uses a retry strategy based on the options provided, such as linear or binary exponential backoff.
// If the maximum number of retries is reached and the function still fails, it returns an error with the provided failure message.
// Example:
//
//	func main() {
//		_, err := RunWithRetry(RetryOptions{
//			RetryCount:  3,
//			WaitSeconds: 5,
//			RetryStrategy: "beb",
//		},
//			"Run failed",
//			func(args ...interface{}) (any, error) {
//				fmt.Println("Pretending to run")
//				return nil, errors.New("bogus err")
//			})
//
//		if err != nil {
//			fmt.Println("We had an error -", err.Error())
//		}
//	}
func RunWithRetry(options RetryOptions, failureMsg string, fn RunnerFunc, args ...any) (data any, err error) {
	tries := options.RetryCount - 1 // -1 to go to zero-based -- default of -1 is also handled
	waitSeconds := options.WaitSeconds

	useBeb := strings.ToLower(options.RetryStrategy) == "beb"

	if useBeb {
		waitSeconds = int(math.Ceil(float64(waitSeconds) / 2.0)) // good starting point for beb
	}

	for {
		if useBeb {
			waitSeconds *= 2
		}

		waitDuration := time.Duration(waitSeconds) * time.Second

		if data, err = fn(args...); err != nil {
			if tries <= 0 {
				return data, errors.New(failureMsg + " - " + err.Error()) // after n retries - maybe we should wrap here
			} else {
				tries--
				fmt.Println("Sleeping for:", waitDuration.String())
				time.Sleep(waitDuration)
				continue
			}
		}

		break
	}
	return
}
