package mmdbserver

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidEditionID  = &Error{StatusCode: http.StatusBadRequest, Message: "Invalid Edition ID"}
	ErrInvalidAccountID  = &Error{StatusCode: http.StatusBadRequest, Message: "Invalid Account ID"}
	ErrInvalidSuffix     = &Error{StatusCode: http.StatusBadRequest, Message: "Invalid Suffix Name"}
	ErrInvalidDate       = &Error{StatusCode: http.StatusBadRequest, Message: "Invalid Date"}
	ErrInvalidMD5Hash    = &Error{StatusCode: http.StatusBadRequest, Message: "Invalid MD5 Hash"}
	ErrInvalidLicenseKey = &Error{StatusCode: http.StatusUnauthorized, Message: "Invalid License Key"}
	ErrMethodNotAllowed  = &Error{StatusCode: http.StatusMethodNotAllowed}
	ErrDatabaseNotFound  = &Error{StatusCode: http.StatusNotFound, Message: "Database not found"}
	ErrDatabaseLatest    = &Error{StatusCode: http.StatusNotModified, Message: "Database is latest"}
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

func WriteError(rw http.ResponseWriter, err error) {
	var exception *Error
	if errors.As(err, &exception) {
		http.Error(rw, err.Error(), exception.StatusCode)
	} else {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	return
}
