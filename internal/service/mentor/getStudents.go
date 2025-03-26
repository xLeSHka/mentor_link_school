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

func (s *MentorService) GetStudents(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.Pair, int64, error) {
	students, total, err := s.mentorRepository.GetStudents(ctx, userID, groupID, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.Pair{}, 0, nil
		}
		return nil, 0, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return students, total, nil
}
