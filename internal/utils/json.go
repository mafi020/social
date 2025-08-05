package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	// whitelist the data
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func JSONErrorResponse(w http.ResponseWriter, status int, message any) error {
	type errorStruct struct {
		Success bool `json:"success"`
		Status  int  `json:"status"`
		Error   any  `json:"error"`
	}
	return WriteJSON(w, status, &errorStruct{Success: false, Status: status, Error: message})
}

func JSONResponse(w http.ResponseWriter, status int, data any) error {
	type responseStruct struct {
		Success bool `json:"success"`
		Status  int  `json:"status"`
		Data    any  `json:"data"`
	}
	return WriteJSON(w, status, &responseStruct{Success: true, Status: status, Data: data})
}
