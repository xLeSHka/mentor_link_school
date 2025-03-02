package mentorService

import (
	"context"
	"errors"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *MentorService) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	own, err := s.mentorRepository.CheckRequest(ctx, request.ID, request.MentorID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !own {
		return httpError.New(http.StatusBadRequest, "Request not found")
	}
	err = s.mentorRepository.UpdateRequest(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpError.New(http.StatusNotFound, err.Error())
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
