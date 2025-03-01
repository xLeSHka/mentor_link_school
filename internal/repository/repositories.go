package repository

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	GetAllGroups(ctx context.Context) ([]*models.Group, error)
}
type UsersRepository interface {
	Create(ctx context.Context, person *models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	GetByEMail(ctx context.Context, email string) (person *models.User, err error)
	EditUser(ctx context.Context, userID uuid.UUID, updates map[string]any) (*models.User, error)
}

type MinioRepository interface {
	UploadImage(file *models.File) (string, error)
	GetImage(image string) (string, error)
	DeleteImage(personID uuid.UUID) error
}
