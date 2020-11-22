package repositories

import "errors"

var (
	// The requested record was not found
	ErrorRecordNotFound = errors.New("record not found")
	// An error occurred while working with the database
	ErrorDatabase = errors.New("database error occurred")

	// Specific to the User model
	// If the username already exists. Can occur when creating a user.
	ErrorUserUsernameAlreadyExists = errors.New("username already exists")
)
