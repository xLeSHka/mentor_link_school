package mentorService

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *MentorService) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	return s.mentorRepository.UpdateRequest(ctx, request)
}
