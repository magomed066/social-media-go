package usersErrors

import "errors"

var ErrPasswordTooLong = errors.New("user password is too long. It must be less than 72")
