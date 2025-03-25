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

func (s *GroupsService) GetMembers(ctx context.Context, groupID uuid.UUID, page, size int) ([]*models.User, int64, error) {
	members, total, err := s.groupRepository.GetMembers(ctx, groupID, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.User{}, 0, nil
		}
		return nil, 0, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return members, total, nil
}
