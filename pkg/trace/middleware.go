package trace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var BackendLogURL string
var APIKey string

// LogEntry represents the structure of a log sent to the backend
type LogEntry struct {
	TraceID     string `json:"trace_id"`
	Method      string `json:"method"`
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	ClientIP    string `json:"client_ip"`
	RequestBody string `json:"request_body,omitempty"`
	ResponseBody string `json:"response_body,omitempty"`
	Latency     string `json:"latency"`
	Timestamp   string `json:"timestamp"`
	Level       LogLevel `json:"level"`
}

// CustomResponseWriter wraps gin.ResponseWriter to capture the response body
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	
	w.body.Write(b) // Capture the response body
	return w.ResponseWriter.Write(b)
}

func InitLogger(backendURL string, apiKey string) {
	BackendLogURL = backendURL
	APIKey = apiKey
}

// GinTraceMiddleware sets up the Trace ID and logs requests and responses
func GinTraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or generate Trace ID
		traceID := c.GetHeader(TraceIDHeader)
		if traceID == "" {
			traceID = GenerateTraceID()
		}
		c.Set(string(TraceIDKey), traceID)
		c.Writer.Header().Set(TraceIDHeader, traceID)

		// Capture request body
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body
		}

		// Start timer
		start := time.Now()

		// Use custom response writer
		customWriter := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = customWriter

		// Process the request
		c.Next()

		// Capture response body
		responseBody := customWriter.body.String()

		// Calculate latency
		latency := time.Since(start)

		// Create log entry
		logEntry := LogEntry{
			TraceID:      traceID,
			Method:       c.Request.Method,
			URL:          c.Request.URL.String(),
			StatusCode:   c.Writer.Status(),
			ClientIP:     c.ClientIP(),
			RequestBody:  requestBody,
			ResponseBody: responseBody,
			Latency:      latency.String(),
			Timestamp:    time.Now().Format(time.RFC3339),
		}
		fmt.Printf("THIS IS THE LOG ENTRY %v", logEntry)
		// Send log to the backend
		if c.Writer.Status() >= 400 {
			Logger("error", logEntry)
		} else {
			Logger("info", logEntry)
		}
	}
}



// LogToBackend sends a log entry to the backend
func Logger(level LogLevel, entry LogEntry) {
	entry.Level = level
	payload, err := json.Marshal(entry)	
	if err != nil {
		fmt.Printf("Failed to marshal log entry: %v", err)
		return
	}
	fmt.Printf("THIS IS THE PAYLOAD %v", payload)
	
	req, err := http.NewRequest("POST", BackendLogURL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("X-Moolya-API-Key", APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send log entry to backend: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Backend log entry returned status code: %d", resp.StatusCode)
	}
}
