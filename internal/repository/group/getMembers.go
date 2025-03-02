package group

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *GroupRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]*models.Role, error) {
	return nil, nil
}
