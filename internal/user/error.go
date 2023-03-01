package user

import (
	"errors"
	"fmt"
)

type ErrNotFound struct {
	UserID string
}

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user '%s' not found", e.UserID)
}
