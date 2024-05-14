package main

import (
	"errors"
	"testing"
	"time"
)

// TestRunWithRetry is a WIP at the moment
func TestRunWithRetry(t *testing.T) {
	var RunFunc RunnerFunc = func(args ...interface{}) (interface{}, error) {
		strike, _ := args[0].(int)
		if strike != 5 {
			return nil, errors.New("return an error")
		}
		return "success", nil
	}

	tests := []struct {
		name           string
		options        RetryOptions
		funcToRun      RunnerFunc
		args           []interface{}
		expectedError  error
		expectedResult interface{}
	}{
		{
			name: "BEB strategy and function succeeds after few runs",
			options: RetryOptions{
				RetryCount:    3,
				WaitSeconds:   1,
				RetryStrategy: "beb",
			},
			funcToRun:      RunFunc,
			args:           []interface{}{5},
			expectedError:  nil,
			expectedResult: "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunWithRetry(tt.options, "Failure after retries", tt.funcToRun, tt.args...)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) ||
				(err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("RunWithRetry() error = %v, expected %v", err, tt.expectedError)
				return
			}
			if result != tt.expectedResult {
				t.Errorf("RunWithRetry() = %v, expected %v", result, tt.expectedResult)
			}
		})
		// waiting for goroutines to finish
		time.Sleep(10 * time.Second)
	}
}
