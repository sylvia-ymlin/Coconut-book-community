package user

import (
	"net/http"

	"github.com/sylvia-ymlin/Coconut-book-community/internal/app"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/response"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/services"
	"github.com/gin-gonic/gin"
)

// RegisterUserDTO 注册用户的请求参数
type RegisterUserDTO struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

const (
	RegisterUserDTO_Username = "username"
	RegisterUserDTO_Password = "password"
)

// UserRegisterHandler registers a new user
// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags User
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} response.RegisterResponse
// @Failure 400 {object} response.CommonResponse
// @Router /user/register/ [post]
func UserRegisterHandler(c *gin.Context) {
	var registerRequest RegisterUserDTO
	var res response.RegisterResponse
	// if err := c.ShouldBindJSON(&registerRequest); err != nil {
	// 	logrus.Debug("UserRegisterHandler error: ", err)
	// 	res.CommonResponse.StatusCode = response.Failed
	// 	res.CommonResponse.StatusMsg = response.ErrInvalidParams
	// 	c.JSON(http.StatusOK, res)
	// 	return
	// }
	var ok1, ok2 bool
	registerRequest.Username, ok1 = c.GetQuery(RegisterUserDTO_Username)
	registerRequest.Password, ok2 = c.GetQuery(RegisterUserDTO_Password)
	if !ok1 || !ok2 {
		res.CommonResponse.StatusCode = response.Failed
		res.CommonResponse.StatusMsg = response.ErrInvalidParams
		c.JSON(http.StatusOK, res)
		return
	}

	res, err := services.CreateUser(registerRequest.Username, registerRequest.Password)
	if err != nil {
		res.CommonResponse.StatusCode = response.Failed
		switch err.Error() {
		case response.ErrUserExists:
			res.CommonResponse.StatusMsg = response.ErrUserExists
		case response.ErrInvalidPassword:
			res.CommonResponse.StatusMsg = response.ErrInvalidPassword
		case response.ErrInvalidUsername:
			res.CommonResponse.StatusMsg = response.ErrInvalidUsername
		default:
			res.CommonResponse.StatusMsg = response.ErrServerInternal
		}
		c.JSON(http.StatusOK, res)
		return
	}
	res.CommonResponse.StatusCode = response.Success
	app.ZeroCheck(res.UserID)
	c.JSON(http.StatusOK, res)
}
