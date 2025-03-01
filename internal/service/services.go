package service

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

type GroupService interface {
	GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error)
	GetGroup(ctx context.Context, mentor *models.Group) (*models.Group, error)
	CreateGroup(ctx context.Context, group *models.Group) error
}
type MentorService interface {
	GetMentorRequest(ctx context.Context, req *models.GetMentorRequest) error
	CreateMentorRequest(ctx context.Context, mentor *models.CreateMentorRequest) error
	GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.Mentor, error)
	GetMentor(ctx context.Context, mentor *models.Mentor) (*models.Mentor, error)
}
type UserService interface {
	Create(ctx context.Context, user *models.User) (string, error)
	Login(ctx context.Context, email string, password string) (*models.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError)
}
