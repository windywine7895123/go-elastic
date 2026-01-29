package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloHandler_Success(t *testing.T) {
	// Setup
	app := fiber.New()
	app.Get("/v2/hello", HelloHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/v2/hello", nil)

	// Execute
	resp, err := app.Test(req, -1)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, result, "message")
	assert.Equal(t, "hello from service", result["message"])
}

func TestHelloHandler_ContentType(t *testing.T) {
	// Setup
	app := fiber.New()
	app.Get("/v2/hello", HelloHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/v2/hello", nil)

	// Execute
	resp, err := app.Test(req, -1)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func BenchmarkHelloHandler(b *testing.B) {
	// Setup
	app := fiber.New()
	app.Get("/v2/hello", HelloHandler)

	req := httptest.NewRequest("GET", "/v2/hello", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req, -1)
	}
}
