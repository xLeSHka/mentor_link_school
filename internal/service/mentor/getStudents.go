package mentorService

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (s *MentorService) GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.PairWithGIDs, error) {
	exist, err := s.usersRepository.CheckExists(ctx, userID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return nil, httpError.New(http.StatusNotFound, "User Not Found")
	}
	students, err := s.mentorRepository.GetStudents(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.PairWithGIDs{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return students, nil
}
