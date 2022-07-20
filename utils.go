package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// getUUID returns a new UUID.
func getUUID() string {
	return uuid.New().String()[:8]
}

// writeData writes the data object to the response writer.
func writeData(w http.ResponseWriter, data Response, code int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data.Code = code
	switch code {
	case http.StatusOK:
		data.Status = "OK"
	case http.StatusForbidden:
		data.Status = "Forbidden"
	case http.StatusUnauthorized:
		data.Status = "Unauthorized"
	}

	text, _ := json.Marshal(data)
	w.Write([]byte(string(text) + "\n"))
}

// getTime returns the current time for timing functions and code blocks
func getTime() time.Time {
	return time.Now()
}

// endTimer logs the elapsed time of a function or code block
func endTimer(startTime time.Time, logMessage string) {
	elapsedTime := time.Since(startTime)
	log.Info(logMessage, zap.Duration("elapsed", elapsedTime))
}
