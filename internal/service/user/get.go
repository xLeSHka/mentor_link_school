package userService

import (
	"context"
	"errors"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"

	"github.com/google/uuid"
)

func (s *UsersService) GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error) {
	user, err := s.usersRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, err.Error())
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}
