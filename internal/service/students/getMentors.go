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

func (s *StudentService) GetMentors(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Role, error) {
	mentors, err := s.studentRepository.GetMentors(ctx, userID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.Role{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return mentors, nil
}
