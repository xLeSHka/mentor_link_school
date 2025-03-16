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

func (s *GroupsService) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {

	group, err := s.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, err.Error())
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return group, nil
}
