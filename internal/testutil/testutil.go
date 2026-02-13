package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// SetupTestRouter creates a test Gin router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// MakeRequest makes an HTTP request for testing
func MakeRequest(router *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// AssertJSON asserts that response body matches expected JSON
func AssertJSON(t *testing.T, w *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(t, err)

	expectedBytes, _ := json.Marshal(expected)
	var expectedMap interface{}
	json.Unmarshal(expectedBytes, &expectedMap)

	assert.Equal(t, expectedMap, actual)
}

// AssertStatusCode asserts HTTP status code
func AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	assert.Equal(t, expectedCode, w.Code)
}

// MockJWTToken generates a mock JWT token for testing
func MockJWTToken(userID uint) string {
	// This is a simplified mock token for testing
	// In real tests, you would generate a proper JWT
	return "mock_jwt_token_" + string(rune(userID))
}
