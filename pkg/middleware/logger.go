package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"mime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"alfredo/tabunganku/pkg/log"
)

// RequestIDKey is the context key for the request ID
const RequestIDKey = "request_id"

// Logger creates a custom logger middleware with request ID correlation
func Logger(logger log.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Get or generate request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}

		// Store in context
		ctx := context.WithValue(c.UserContext(), RequestIDKey, requestID)
		c.SetUserContext(ctx)
		c.Locals(RequestIDKey, requestID)

		// Process request
		err := c.Next()

		// Calculate response time
		responseTime := time.Since(start)

		// Get context-aware logger
		ctxLogger := logger.WithContext(ctx)

		// Create log arguments
		logArgs := []any{
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"ip", c.IP(),
			"latency", responseTime.String(),
			"user_agent", c.Get("User-Agent"),
		}

		// Include POST data if it's a POST request and not binary data
		if c.Method() == "POST" && !isBinaryContent(c.Get("Content-Type")) {
			body := c.Body()
			if len(body) > 0 && isValidJSON(body) {
				// Pretty-print the JSON for better readability
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
					// If the body is too large, truncate it
					bodyStr := prettyJSON.String()
					if len(bodyStr) > 1000 {
						bodyStr = bodyStr[:1000] + "... [truncated]"
					}

					var bodyMap map[string]interface{}
					if err := json.Unmarshal(body, &bodyMap); err == nil {
						if _, ok := bodyMap["password"]; ok {
							bodyMap["password"] = "***"
						}

						if _, ok := bodyMap["refresh_token"]; ok {
							bodyMap["refresh_token"] = "***"
						}

						if maskedBody, err := json.MarshalIndent(bodyMap, "", "  "); err == nil {
							logArgs = append(logArgs, "body", string(maskedBody))
						} else {
							logArgs = append(logArgs, "body", string(body))
						}
					} else {
						logArgs = append(logArgs, "body", string(body))
					}

				} else {
					// If not valid JSON, include raw body (truncated if needed)
					bodyStr := string(body)
					if len(bodyStr) > 1000 {
						bodyStr = bodyStr[:1000] + "... [truncated]"
					}
					logArgs = append(logArgs, "body", bodyStr)
				}
			}
		}

		// Log the request details
		ctxLogger.Info("HTTP Request", logArgs...)

		return err
	}
}

// isBinaryContent checks if the content type is likely binary data
func isBinaryContent(contentType string) bool {
	if contentType == "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	// Check for common binary content types
	return strings.HasPrefix(mediaType, "image/") ||
		strings.HasPrefix(mediaType, "audio/") ||
		strings.HasPrefix(mediaType, "video/") ||
		strings.HasPrefix(mediaType, "application/octet-stream") ||
		strings.HasPrefix(mediaType, "application/pdf") ||
		strings.HasPrefix(mediaType, "application/zip")
}

// isValidJSON checks if a byte slice contains valid JSON
func isValidJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
