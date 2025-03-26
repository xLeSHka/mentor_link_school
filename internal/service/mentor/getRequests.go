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

func (s *MentorService) GetMyHelps(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.HelpRequest, int64, error) {
	helps, total, err := s.mentorRepository.GetMyHelpers(ctx, userID, groupID, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.HelpRequest{}, 0, nil
		}
		return nil, 0, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return helps, total, nil
}
