package usersErrors

import "errors"

var ErrPasswordTooLong = errors.New("user password is too long. It must be less than 72")

var ErrDuplicateEmail = errors.New("Email already exists")
var ErrDuplicateUsername = errors.New("Username already exists")
