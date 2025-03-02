package groupService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *GroupsService) GetMembers(ctx context.Context, ownerID, groupId uuid.UUID) ([]*models.Role, error) {
	exist, err := s.groupRepository.CheckGroupExists(ctx, ownerID, groupId)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return nil, httpError.New(http.StatusNotFound, "Group not found")
	}
	members, err := s.groupRepository.GetMembers(ctx, groupId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.Role{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return members, nil
}
