package laundry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LaundryError struct {
	Reasons []string `json:"errors"`
	Status  int      `json:"status"`
}

func ExtError(e error) *LaundryError {
	return &LaundryError{[]string{e.Error()}, http.StatusBadRequest}
}

func NewError(e string) *LaundryError {
	return &LaundryError{[]string{e}, http.StatusBadRequest}
}

func (e LaundryError) Error() string {
	return fmt.Sprintf("Error (%d): %s", e.Status, strings.Join(e.Reasons, ", "))
}

func (e *LaundryError) WithStatus(i int) *LaundryError {
	e.Status = i

	return e
}

func (e *LaundryError) Add(r string) *LaundryError {
	e.Reasons = append(e.Reasons, r)

	return e
}

func (e *LaundryError) AsJSON() []byte {
	j, _ := json.Marshal(e)

	return j
}