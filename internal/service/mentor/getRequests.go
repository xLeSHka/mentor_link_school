package mentorService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *MentorService) GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error) {
	helps, err := s.usersRepository.GetMyRequests(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, err.Error())
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return helps, nil
}
