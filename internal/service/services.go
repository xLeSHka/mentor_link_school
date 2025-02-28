package service

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) (string, error)
	Login(ctx context.Context, email string, password string) (*models.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError)
}
