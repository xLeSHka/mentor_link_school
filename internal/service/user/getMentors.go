package userService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *UsersService) GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.RoleWithGIDs, error) {
	exist, err := s.usersRepository.CheckExists(ctx, userID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return nil, httpError.New(http.StatusNotFound, "User Not Found")
	}
	mentors, err := s.usersRepository.GetMentors(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.RoleWithGIDs{}, nil
		}
		return nil, err
	}
	return mentors, nil
}
