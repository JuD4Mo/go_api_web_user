package user

import "errors"

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")

var ErrUserNotFound = errors.New("user does not exist")
