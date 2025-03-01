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
	GetMentorRequest(ctx context.Context, data *models.HelpRequest) error
	CreateMentorRequest(ctx context.Context, data *models.CreateMentorRequest) error
	GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.Mentor, error)
	GetMentor(ctx context.Context, mentor *models.Mentor) (*models.Mentor, error)
}
type UsersRepository interface {
	Login(ctx context.Context, person *models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	EditUser(ctx context.Context, userID uuid.UUID, updates map[string]any) (*models.User, error)
}

type MinioRepository interface {
	UploadImage(file *models.File) (string, error)
	GetImage(image string) (string, error)
	DeleteImage(personID uuid.UUID) error
}
