package carefulness

import "errors"

var ErrEmailAlreadyExists = errors.New("user with this email already exists")
var InvalidCreditantials = errors.New("Invalid credetantials")
var UserNotFound = errors.New("User not found")
