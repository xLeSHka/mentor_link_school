package usersService

import (
	"context"
	"errors"
	"github.com/xLeSHka/mentorLinkSchool/internal/utils/password"
	"gorm.io/gorm"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (s *UserService) Login(ctx context.Context, telegram, pw string) (string, error) {

	user, err := s.usersRepository.Login(ctx, telegram, pw)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", httpError.New(http.StatusNotFound, "User not found")
		}
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	err = password.Compare([]byte(pw), user.Password, s.cryptoKey)
	if err != nil {
		return "", httpError.New(http.StatusUnauthorized, "Invalid password")
	}
	token, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	err = s.cache.SaveToken(ctx, user.ID, token)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}
