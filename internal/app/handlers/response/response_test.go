package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommonResponse(t *testing.T) {
	t.Run("success response", func(t *testing.T) {
		resp := CommonResponse{
			StatusCode: Success,
			StatusMsg:  "操作成功",
		}

		assert.Equal(t, int32(Success), resp.StatusCode)
		assert.Equal(t, "操作成功", resp.StatusMsg)
	})

	t.Run("failed response", func(t *testing.T) {
		resp := CommonResponse{
			StatusCode: Failed,
			StatusMsg:  ErrServerInternal,
		}

		assert.Equal(t, int32(Failed), resp.StatusCode)
		assert.Equal(t, ErrServerInternal, resp.StatusMsg)
	})
}

func TestRegisterResponse(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		resp := RegisterResponse{
			CommonResponse: CommonResponse{
				StatusCode: Success,
				StatusMsg:  "注册成功",
			},
			UserID: 123,
			Token:  "test_token_abc123",
		}

		assert.Equal(t, int32(Success), resp.StatusCode)
		assert.Equal(t, 123, resp.UserID)
		assert.Equal(t, "test_token_abc123", resp.Token)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("failed registration", func(t *testing.T) {
		resp := RegisterResponse{
			CommonResponse: CommonResponse{
				StatusCode: Failed,
				StatusMsg:  ErrUserExists,
			},
		}

		assert.Equal(t, int32(Failed), resp.StatusCode)
		assert.Equal(t, ErrUserExists, resp.StatusMsg)
		assert.Equal(t, 0, resp.UserID)
		assert.Empty(t, resp.Token)
	})
}

func TestGetUserInfoResponse(t *testing.T) {
	t.Run("successful user info", func(t *testing.T) {
		resp := GetUserInfoResponse{
			CommonResponse: CommonResponse{
				StatusCode: Success,
			},
			User: User{
				ID:            1,
				Name:          "testuser",
				FollowCount:   10,
				FollowerCount: 20,
				IsFollow:      true,
			},
		}

		assert.Equal(t, int32(Success), resp.StatusCode)
		assert.Equal(t, 1, resp.User.ID)
		assert.Equal(t, "testuser", resp.User.Name)
		assert.Equal(t, 10, resp.User.FollowCount)
		assert.Equal(t, 20, resp.User.FollowerCount)
		assert.True(t, resp.User.IsFollow)
	})
}

func TestErrorMessages(t *testing.T) {
	t.Run("error constants", func(t *testing.T) {
		assert.NotEmpty(t, ErrInvalidParams)
		assert.NotEmpty(t, ErrUserExists)
		assert.NotEmpty(t, ErrUserNotExists)
		assert.NotEmpty(t, ErrUserPassword)
		assert.NotEmpty(t, ErrServerInternal)
		assert.NotEmpty(t, ErrUserToken)
		assert.NotEmpty(t, ErrInvalidUsername)
		assert.NotEmpty(t, ErrInvalidPassword)

		// Verify they are different
		assert.NotEqual(t, ErrUserExists, ErrUserNotExists)
		assert.NotEqual(t, ErrInvalidParams, ErrServerInternal)
	})
}

func TestStatusCodes(t *testing.T) {
	t.Run("status code constants", func(t *testing.T) {
		assert.Equal(t, 0, Success)
		assert.Equal(t, 500, Failed)
		assert.Equal(t, 401, TokenExpired)
		assert.NotEqual(t, Success, Failed)
	})
}
