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

func (r *UsersService) GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Role, error) {
	gr, err := r.usersRepository.GetGroups(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return []*models.Role{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return gr, nil
}
