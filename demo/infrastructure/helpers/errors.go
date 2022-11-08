package helpers

import (
	"errors"
)

func FormatError(err string) error {
	return errors.New("Error")
}

var ErrNotFound = errors.New("Not found")

var ErrInvalidEntity = errors.New("Invalid entity")

var ErrCannotBeDeleted = errors.New("Cannot Be Deleted")
