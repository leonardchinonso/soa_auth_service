package errors

import (
	"errors"
	"net/http"
)

const (
	// ErrInvalidLogin for when login credentials are incorrect
	ErrInvalidLogin = "invalid login credentials"
)

// RestError is the custom struct for a request error
type RestError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Err     string      `json:"error"`
	Data    interface{} `json:"data"`
}

// Error returns the message from RestError
// it fulfills the interface requirements for the standard error type
func (re *RestError) Error() string {
	return re.Message
}

// Status checks if the error is of the RestError type
// returns an internal server error if it is not of the set type
func Status(err error) int {
	var re *RestError
	if errors.As(err, &re) {
		return re.Status
	}
	return http.StatusInternalServerError
}

// ErrBadRequest returns a RestError for a bad request
func ErrBadRequest(message string, data interface{}) *RestError {
	return &RestError{
		Status:  http.StatusBadRequest,
		Message: message,
		Err:     "Bad Request",
		Data:    data,
	}
}

// ErrInternalServerError returns a RestError for internal server error
func ErrInternalServerError(message string, data interface{}) *RestError {
	return &RestError{
		Status:  http.StatusInternalServerError,
		Message: message,
		Err:     "Internal Server Error",
		Data:    data,
	}
}

// ErrUnauthorized returns a RestError for an unauthorized request
func ErrUnauthorized(message string, data interface{}) *RestError {
	return &RestError{
		Status:  http.StatusUnauthorized,
		Message: message,
		Err:     "Unauthorized",
		Data:    data,
	}
}

// ErrorToStringSlice converts a slice of errors to a slice of string
func ErrorToStringSlice(errs []error) []string {
	var errStrings []string
	for _, e := range errs {
		errStrings = append(errStrings, e.Error())
	}
	return errStrings
}
