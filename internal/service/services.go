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
	GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	UpdateRequest(ctx context.Context, request *models.HelpRequest) error
	GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
}
type UserService interface {
	Login(ctx context.Context, name string) (*models.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError)
	GetMyMentors(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
	GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	CreateRequest(ctx context.Context, request *models.HelpRequest) error
	GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
}
