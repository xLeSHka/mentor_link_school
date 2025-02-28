package userService

import (
	"context"
	"errors"
	"net/http"
	"prodapp/internal/app/helpers"
	"prodapp/internal/app/httpError"
	"prodapp/internal/models"

	"gorm.io/gorm"
)

func (s *UsersService) Create(ctx context.Context, user *models.User) (string, error) {
	hashPassword, _ := helpers.Encrypt(user.Password, s.cryptoKey)
	user.Password = hashPassword

	_, err := s.usersRepository.Create(ctx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", httpError.New(http.StatusConflict, "Такой email уже зарегистрирован.")
		}
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}

	token, err := s.GenerateAccessToken(user.ID)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}
