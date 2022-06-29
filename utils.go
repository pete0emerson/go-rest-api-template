package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// getUUID returns a new UUID.
func getUUID() string {
	return uuid.New().String()[:8]
}

// writeData writes the data object to the response writer.
func writeData(w http.ResponseWriter, data map[string]string, code int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data["code"] = strconv.Itoa(code)
	switch code {
	case http.StatusOK:
		data["status"] = "OK"
	case http.StatusForbidden:
		data["status"] = "Forbidden"
	case http.StatusUnauthorized:
		data["status"] = "Unauthorized"
	}

	text, _ := json.Marshal(data)
	w.Write([]byte(string(text) + "\n"))
}
