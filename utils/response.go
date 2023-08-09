package utils

import "net/http"

// Response formats the standard http response
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseStatusOK returns a Response with StatusOK
func ResponseStatusOK(message string, data interface{}) *Response {
	return &Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

// ResponseStatusCreated returns a Response with StatusCreated
func ResponseStatusCreated(message string, data interface{}) *Response {
	return &Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}
