package app

import (
	"github.com/sylvia-ymlin/Coconut-book-community/config"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type User struct {
	models.UserModel
	Token  string            `json:"token"`
	Claims *jwt.CustomClaims `json:"claims"`
}

const UserKeyName = "user"

func ZeroCheck[T comparable](v ...T) bool {
	if !config.IsDebug() {
		return false
	}
	var zero T
	for _, item := range v {
		if item == zero {
			logrus.Errorf("zero value: %v", v)
		}
	}
	return false
}
