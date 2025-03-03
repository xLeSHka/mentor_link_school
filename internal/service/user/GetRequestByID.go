package userService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (r *UsersService) GetRequestByID(ctx context.Context, reqID uuid.UUID) (models.HelpRequest, error) {
	req, err := r.usersRepository.GetRequestByID(ctx, reqID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.HelpRequest{}, httpError.New(http.StatusNotFound, err.Error())
		}
		return models.HelpRequest{}, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return req, err
}
