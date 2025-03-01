package repository

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error)
	GetGroup(ctx context.Context, group *models.Group) (*models.Group, error)
}
type MentorRepository interface {
	GetMyHelpers(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	UpdateRequest(ctx context.Context, request *models.HelpRequest) error
	GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
}
type UsersRepository interface {
	Login(ctx context.Context, person *models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	EditUser(ctx context.Context, userID uuid.UUID, updates map[string]any) (*models.User, error)
	GetMyMentors(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
	GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
	GetMyRequests(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	CreateRequest(ctx context.Context, request *models.HelpRequest) error
}

type MinioRepository interface {
	UploadImage(file *models.File) (string, error)
	GetImage(image string) (string, error)
	DeleteImage(personID uuid.UUID) error
}
