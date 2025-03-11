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

func (s *UsersService) GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequestWithGIDs, error) {
	exist, err := s.usersRepository.CheckExists(ctx, userID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return nil, httpError.New(http.StatusNotFound, "User Not Found")
	}
	requests, err := s.usersRepository.GetMyRequests(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.HelpRequestWithGIDs{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return requests, nil
}
