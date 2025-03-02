package groupService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *GroupsService) GetMembers(ctx context.Context, groupId uuid.UUID) ([]*models.Role, error) {
	return nil, nil
}
