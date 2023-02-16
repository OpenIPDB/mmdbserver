package mmdbserver

import (
	"fmt"
	"net/http"
)

var (
	ErrInvalidEditionID = &Error{StatusCode: http.StatusNotFound, Message: "Invalid Edition ID"}
	ErrInvalidMD5Hash   = &Error{StatusCode: http.StatusNotFound, Message: "Invalid MD5 Hash"}
	ErrInvalidAccount   = &Error{StatusCode: http.StatusBadRequest, Message: "Invalid Account"}
	ErrInvalidRequest   = &Error{StatusCode: http.StatusBadRequest, Message: "Invalid Request"}
	ErrUnauthorized     = &Error{StatusCode: http.StatusUnauthorized}
	ErrMethodNotAllowed = &Error{StatusCode: http.StatusMethodNotAllowed}
	ErrDatabaseNotFound = &Error{StatusCode: http.StatusNotFound, Message: "Database not found"}
	ErrDatabaseLatest   = &Error{StatusCode: http.StatusNotModified, Message: "Database is latest"}
)

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	message := e.Message
	if message == "" {
		message = http.StatusText(e.StatusCode)
	}
	return fmt.Sprintf("%d %s", e.StatusCode, message)
}
