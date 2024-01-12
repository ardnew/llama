package walk

import "fmt"

type (
	RunError struct {
		error
		src error
	}
)

func newRunError(err error) *RunError {
	if err == nil {
		return nil
	}
	return &RunError{
		error: fmt.Errorf("runtime error: %w", err),
		src:   err,
	}
}

func (e *RunError) Error() string {
	if e == nil || e.src == nil {
		return ""
	}
	return e.error.Error()
}

func (e *RunError) Unwrap() []error {
	if e == nil || e.src == nil {
		return nil
	}
	return []error{e.src}
}
