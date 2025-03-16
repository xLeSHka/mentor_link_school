package usersService

import (
	"context"
	"errors"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *UserService) Register(ctx context.Context, user *models.User) (string, error) {
	user, err := s.usersRepository.Register(ctx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", httpError.New(http.StatusConflict, err.Error())
		}
		return "", httpError.New(http.StatusInternalServerError, err.Error())
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
