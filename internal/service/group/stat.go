package groupService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *GroupsService) GetStat(ctx context.Context, groupID uuid.UUID) (*models.GroupStat, error) {
	return nil, nil
}
