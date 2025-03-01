package mentorService

import (
	"context"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *MentorService) GetMentor(ctx context.Context, mentor *models.Mentor) (*models.Mentor, error) {
	return s.mentorRepository.GetMentor(ctx, mentor)
}

func (s *MentorService) GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.Mentor, error) {
	return s.mentorRepository.GetMentors(ctx, userID)
}
