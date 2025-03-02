package userService

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *UsersService) GetGroups(ctx context.Context, userID uuid.UUID, role string) ([]*models.Group, error) {
	gr, err := r.usersRepository.GetGroups(ctx, userID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return []*models.Group{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return gr, nil
}
