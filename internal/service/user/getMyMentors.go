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

func (s *UsersService) GetMyMentors(ctx context.Context, userID uuid.UUID) ([]*models.PairWithGIDs, error) {
	exist, err := s.usersRepository.CheckExists(ctx, userID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return nil, httpError.New(http.StatusNotFound, "User Not Found")
	}
	mentors, err := s.usersRepository.GetMyMentors(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.PairWithGIDs{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return mentors, nil
}
