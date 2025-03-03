package groupService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"

	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (s *GroupsService) Create(ctx context.Context, group *models.Group, userID uuid.UUID) error {

	err := s.groupRepository.Create(ctx, group, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return httpError.New(http.StatusConflict, "Такая организация уже зарегистрирована.")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
