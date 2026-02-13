package recommendation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestSearchBooksHandler(t *testing.T) {
	router := setupTestRouter()
	router.GET("/search", SearchBooksHandler)

	t.Run("valid search query", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=golang&top_k=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		statusCode, ok := response["status_code"].(float64)
		assert.True(t, ok)
		assert.Equal(t, float64(0), statusCode)

		books, ok := response["books"].([]interface{})
		assert.True(t, ok)
		assert.Greater(t, len(books), 0)
		assert.LessOrEqual(t, len(books), 5)
	})

	t.Run("search with chinese query", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=计算机&top_k=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		statusCode := response["status_code"].(float64)
		assert.Equal(t, float64(0), statusCode)
	})

	t.Run("missing query parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?top_k=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		statusCode := response["status_code"].(float64)
		assert.NotEqual(t, float64(0), statusCode)
	})

	t.Run("default top_k", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		books, ok := response["books"].([]interface{})
		assert.True(t, ok)
		assert.LessOrEqual(t, len(books), 10) // Default top_k is 10
	})

	t.Run("custom top_k", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=programming&top_k=3", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		books, ok := response["books"].([]interface{})
		assert.True(t, ok)
		assert.LessOrEqual(t, len(books), 3)
	})
}

func TestGetRecommendationsHandler_Structure(t *testing.T) {
	t.Run("handler function exists", func(t *testing.T) {
		assert.NotNil(t, GetRecommendationsHandler)
	})

	t.Run("search handler function exists", func(t *testing.T) {
		assert.NotNil(t, SearchBooksHandler)
	})
}

func TestRecommendationService(t *testing.T) {
	t.Run("service is initialized", func(t *testing.T) {
		assert.NotNil(t, recommendService)
	})
}
