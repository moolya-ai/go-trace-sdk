package trace

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var BackendLogURL string

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

func InitLogger(backendURL string) {
	BackendLogURL = backendURL
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
		}

		// Send log to the backend
		Logger("info", logEntry)
	}
}



// LogToBackend sends a log entry to the backend
func Logger(level LogLevel, entry LogEntry) {
	entry.Level = level
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
