package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlanetError struct {
	Code    string
	Message string
	Status  int
}

func (e *PlanetError) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.Message)
}

func (e *PlanetError) JSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.Code,
		Message: e.Message,
	})
}

func ToPlanetError(err error) *PlanetError {
	planetError, ok := err.(*PlanetError)
	if !ok {
		return nil
	}

	return planetError
}

// NewError creates a new error with the given status code and message.
func NewError(statusCode int, message string) error {
	return NewErrorf(statusCode, message)
}

// NewErrorf creates a new error with the given status code and the message
// formatted according to args and format.
func NewErrorf(status int, format string, args ...interface{}) error {
	return &PlanetError{
		Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_"),
		Message: fmt.Sprintf(format, args...),
		Status:  status,
	}
}
