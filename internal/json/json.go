package json

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	setHeaderConfig(w, status)
	return json.NewEncoder(w).Encode(envelope{
		Success: true,
		Data:    data,
	})
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	setHeaderConfig(w, status)

	return json.NewEncoder(w).Encode(envelope{
		Success: true,
		Data:    nil,
		Error:   message,
	})
}

func setHeaderConfig(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
