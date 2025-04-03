package middleware

import (
	"bytes"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"io"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupMiddleware configures middleware for the router
func SetupMiddleware(router *gin.Engine) {
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Request logging middleware
	router.Use(RequestLogger())

	// Recovery middleware
	router.Use(gin.Recovery())
}

// RequestLogger logs request and response details
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read the request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Capture the response
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Log request and response details
		duration := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		userAgent := c.Request.UserAgent()

		// Truncate large request/response bodies to prevent logging too much data
		const maxBodySize = 1024 * 10 // 10KB
		truncateBody := func(body []byte) string {
			if len(body) > maxBodySize {
				return string(body[:maxBodySize]) + "...(truncated)"
			}
			return string(body)
		}

		// Determine log level based on status code
		logLevel := zap.InfoLevel
		if status >= 400 && status < 500 {
			logLevel = zap.WarnLevel
		} else if status >= 500 {
			logLevel = zap.ErrorLevel
		}

		// Don't log large media files and similar content
		contentType := c.GetHeader("Content-Type")
		shouldLogBody := !isMediaContentType(contentType)

		// Create structured log
		fields := []zap.Field{
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("user_agent", userAgent),
		}

		// Only add request/response body for appropriate content types
		if shouldLogBody {
			// Add request body (if not too large or sensitive)
			if len(requestBody) > 0 && !isRequestSensitive(path) {
				fields = append(fields, zap.String("request_body", truncateBody(requestBody)))
			}

			// Add response body (if not too large or sensitive)
			if responseWriter.body.Len() > 0 && !isResponseSensitive(path) {
				fields = append(fields, zap.String("response_body", truncateBody(responseWriter.body.Bytes())))
			}
		}

		// Log with appropriate level
		logger.GetLogger().Log(logLevel, "HTTP Request", fields...)
	}
}

// Helper function to check if a request path contains sensitive information
func isRequestSensitive(path string) bool {
	// Add paths that may contain sensitive information
	sensitivePaths := []string{
		"/login",
		"/register",
		"/users",
		"/auth",
		"/password",
	}

	for _, p := range sensitivePaths {
		if bytes.Contains([]byte(path), []byte(p)) {
			return true
		}
	}
	return false
}

// Helper function to check if a response path contains sensitive information
func isResponseSensitive(path string) bool {
	// Add paths that may return sensitive information
	sensitivePaths := []string{
		"/users",
		"/profile",
		"/auth",
	}

	for _, p := range sensitivePaths {
		if bytes.Contains([]byte(path), []byte(p)) {
			return true
		}
	}
	return false
}

// Helper function to check if content type is media
func isMediaContentType(contentType string) bool {
	mediaContentTypes := []string{
		"image/",
		"video/",
		"audio/",
		"application/pdf",
		"application/zip",
		"application/octet-stream",
	}

	for _, media := range mediaContentTypes {
		if bytes.Contains([]byte(contentType), []byte(media)) {
			return true
		}
	}
	return false
}

// Custom response writer that captures the response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString captures the response body as string
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
