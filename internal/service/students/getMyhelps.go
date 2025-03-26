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

func (s *StudentService) GetMyHelps(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.HelpRequest, int64, error) {
	requests, total, err := s.studentRepository.GetMyRequests(ctx, userID, groupID, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.HelpRequest{}, 0, nil
		}
		return nil, 0, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return requests, total, nil
}
