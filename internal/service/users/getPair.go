package usersService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (s *UserService) GetPair(ctx context.Context, userID, userID2, groupID uuid.UUID) (*models.Pair, error) {
	pair, err := s.usersRepository.GetPair(ctx, userID, userID2, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return pair, nil
}
