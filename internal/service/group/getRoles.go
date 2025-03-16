package groupService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *GroupsService) GetRoles(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Role, error) {
	roles, err := s.groupRepository.GetRoles(ctx, userID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, "user is not have role in group or user is not exist")
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	err = s.cacheRepository.AddRoles(ctx, roles)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return roles, nil
}
