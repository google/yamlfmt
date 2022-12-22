package basic

import "fmt"

type BasicFormatterError struct {
	err error
}

func (e BasicFormatterError) Error() string {
	return fmt.Sprintf("basic formatter error: %v", e.err)
}

func (e BasicFormatterError) Unwrap() error {
	return e.err
}

// func wrapBasicFormatterError(err error) error {
// 	return BasicFormatterError{err: err}
// }
