package studentService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *StudentService) GetMyHelps(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error) {
	requests, err := s.studentRepository.GetMyRequests(ctx, userID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.HelpRequest{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return requests, nil
}
