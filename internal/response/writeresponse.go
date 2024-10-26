package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int `json:"status_code"`
	Data       any `json:"data"`
}

func WriteResponse(w http.ResponseWriter, status int, val any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&Response{
		StatusCode: status,
		Data:       val,
	})
}
