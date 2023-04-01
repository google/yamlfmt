package collections

import "errors"

type Errors []error

func (errs Errors) Combine() error {
	errMessage := ""

	for _, err := range errs {
		if err != nil {
			errMessage += err.Error() + "\n"
		}
	}

	if len(errMessage) == 0 {
		return nil
	}
	return errors.New(errMessage)
}
