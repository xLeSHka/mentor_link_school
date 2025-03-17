package usersService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *UserService) Edit(ctx context.Context, userID uuid.UUID, user *models.User) (*models.User, error) {

	updated, err := s.usersRepository.EditUser(ctx, userID, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, "user not found")
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, httpError.New(http.StatusConflict, "duplicate key")
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return updated, nil
}
