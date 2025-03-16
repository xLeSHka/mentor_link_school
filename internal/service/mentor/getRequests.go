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

func (s *MentorService) GetMyHelps(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error) {
	helps, err := s.mentorRepository.GetMyHelpers(ctx, userID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.HelpRequest{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return helps, nil
}
