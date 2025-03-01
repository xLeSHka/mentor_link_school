package mentorService

import (
	"context"
	"errors"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *GroupsService) CreateGroup(ctx context.Context, group *models.Group) error {

	err := s.groupRepository.Create(ctx, group)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return httpError.New(http.StatusConflict, "Такой email уже зарегистрирован.")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
