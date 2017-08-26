package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  int      `json:"status"`
	Reasons []string `json:"reasons"`
}

func (api *LaundryAPI) RespondJSON(e ErrorResponse, w http.ResponseWriter) {
	// Set response status code
	w.WriteHeader(e.Status)

	j, _ := json.Marshal(e)

	w.Write(j)
}
