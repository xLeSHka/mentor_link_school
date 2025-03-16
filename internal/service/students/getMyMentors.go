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

func (s *StudentService) GetMyMentors(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Pair, error) {
	mentors, err := s.studentRepository.GetMyMentors(ctx, userID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.Pair{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return mentors, nil
}
