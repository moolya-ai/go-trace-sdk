package trace

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
	ErrorLevel LogLevel = "error"
)

// CustomLogEntry represents the structure of a custom log
type CustomLogEntry struct {
	TraceID  string    `json:"trace_id"`
	Level    LogLevel  `json:"level"`
	Message  string    `json:"message"`
	Details  string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// LogWithLevel sends a custom log to the backend
func LogWithLevel(traceID string, level LogLevel, message string, details string) {
	entry := CustomLogEntry{
		TraceID:  traceID,
		Level:    level,
		Message:  message,
		Details:  details,
		Timestamp: time.Now(),
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	resp, err := http.Post(BackendLogURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Failed to send log entry to backend: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Backend log entry returned status code: %d", resp.StatusCode)
	}
}
