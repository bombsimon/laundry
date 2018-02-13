package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// LaundryError represent an error caught in the laundry service.
// The error contains a list of reasons for the error (strings) and
// a status integer to be used as HTTP status code
type LaundryError struct {
	Reasons []string `json:"errors"`
	Status  int      `json:"status"`
	Origin  error    `json:"-"`
}

// New will take a string or an error and create a LaundryError. The status code
// will always default to http.StatusBadRequest
func New(e interface{}) *LaundryError {
	le := LaundryError{
		Status: http.StatusBadRequest,
	}

	switch v := e.(type) {
	case string:
		le.Reasons = []string{v}
	case error:
		le.Reasons = []string{v.Error()}
		le.Origin = v
	}

	return &le
}

// Error makes sure LaundryError implements the error interface and returns
// a string representation of the LaundryError
func (e LaundryError) Error() string {
	return fmt.Sprintf("Error (%d): %s", e.Status, strings.Join(e.Reasons, ", "))
}

// WithStatus will override the default status (http.StatusBadRequest) with
// given integer.
// Usage: err := NewError("Something went wrong").WithStatus(http.StatusNotFound)
func (e *LaundryError) WithStatus(i int) *LaundryError {
	e.Status = i

	return e
}

// CasuedBy will add an error to the LaundryError to see what caused it
func (e *LaundryError) CausedBy(err error) *LaundryError {
	e.Origin = err

	return e
}

// Add will add another error to the existing LaundryError.
// Usage: err := NewError("First error").WithStatus(http.StatusNotFound).Add("Couldn't find it")
func (e *LaundryError) Add(r string) *LaundryError {
	e.Reasons = append(e.Reasons, r)

	return e
}

// AsJSON will marshal the LaundryError to JSON data usable to render from
// a HTTP server
func (e *LaundryError) AsJSON() []byte {
	j, _ := json.Marshal(e)

	return j
}
