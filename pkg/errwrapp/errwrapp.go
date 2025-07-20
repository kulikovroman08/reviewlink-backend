package errwrapp

import (
	"fmt"
	"runtime"
)

func WithCaller(err error) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}

	return fmt.Errorf("%s:%d: %w", file, line, err)
}
