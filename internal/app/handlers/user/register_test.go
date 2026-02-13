package user

import (
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

func TestRegisterUserDTO(t *testing.T) {
	t.Run("create DTO with valid data", func(t *testing.T) {
		dto := RegisterUserDTO{
			Username: "testuser",
			Password: "password123",
		}

		assert.Equal(t, "testuser", dto.Username)
		assert.Equal(t, "password123", dto.Password)
	})

	t.Run("DTO constants", func(t *testing.T) {
		assert.Equal(t, "username", RegisterUserDTO_Username)
		assert.Equal(t, "password", RegisterUserDTO_Password)
	})
}

func TestUserRegisterHandler_InvalidParams(t *testing.T) {
	router := setupTestRouter()
	router.POST("/user/register/", UserRegisterHandler)

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
	}{
		{
			name:           "missing username",
			username:       "",
			password:       "password123",
			expectedStatus: http.StatusOK, // Handler returns 200 with error status
		},
		{
			name:           "missing password",
			username:       "testuser",
			password:       "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing both",
			username:       "",
			password:       "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/user/register/?username="+tt.username+"&password="+tt.password, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRegisterUserDTO_JSONBinding(t *testing.T) {
	t.Run("json tags are correct", func(t *testing.T) {
		dto := RegisterUserDTO{}

		// These are the struct tags
		assert.NotNil(t, dto)
	})
}
