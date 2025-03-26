package usersService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"gorm.io/gorm"
	"net/http"
)

func (s *UserService) GetTelegramID(ctx context.Context, userID uuid.UUID) (int64, error) {
	id, err := s.cache.GetID(ctx, userID)
	if err == nil {
		return id, nil
	}
	user, err := s.usersRepository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, httpError.New(http.StatusNotFound, err.Error())
		}
		return 0, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if user.TelegramID == nil {
		return 0, httpError.New(http.StatusBadRequest, "user fave not telegram id")
	}
	err = s.cache.SaveID(ctx, userID, *user.TelegramID)
	if err != nil {
		return 0, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return *user.TelegramID, nil
}
