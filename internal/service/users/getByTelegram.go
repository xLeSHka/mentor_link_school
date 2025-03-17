package usersService

import (
	"context"
	"errors"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (s *UserService) GetByTelegram(ctx context.Context, telegram string) (*models.User, error) {

	user, err := s.usersRepository.GetByTelegram(ctx, telegram)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, err.Error())
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}
