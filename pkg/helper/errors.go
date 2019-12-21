package helper

import (
	"errors"
)

var (
	// Not Found Error
	ErrNotFound = errors.New("not found")
	// User Not Provided Error
	ErrNoUserProvided = errors.New("user not provided")
	// Database Connection Error
	ErrDBConnection = errors.New("database connection error")
	// Password Not Equal Error
	ErrPasswordsNotEqual = errors.New("password not equal to confirmed password")
	// User Already Exists Error
	ErrUserAlreadyExists = errors.New("user already exists")
	// No User Name Provided Error
	ErrNoUsernameProvided = errors.New("no user name provided")
	// No Password Provided Error
	ErrNoPasswordProvided = errors.New("no password provided")
	// Invalid Username Provided Error
	ErrInvalidUsernameProvided = errors.New("invalid username provided")
	// Invalid Password Provided Error
	ErrInvalidPasswordProvided = errors.New("invalid password provided")
	// Invalid Arguments Provided Error
	ErrInvalidArgumentProvided = errors.New("invalid argument/s provided")
	// Inconsistent IDs Error
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	// Already Exists Error
	ErrAlreadyExists = errors.New("already exists")
	// Inconsistent Mapping Between Route and Handler Error
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)