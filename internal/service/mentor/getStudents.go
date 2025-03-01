package mentorService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *MentorService) GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error) {
	return s.mentorRepository.GetStudents(ctx, userID)
}
