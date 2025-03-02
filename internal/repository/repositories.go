package repository

import (
	"context"

	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group, userID uuid.UUID) error
	UpdateToMentor(ctx context.Context, groupID, userID uuid.UUID) error
	UpdateInviteCode(ctx context.Context, groupID uuid.UUID, inviteCode string) error
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]*models.Role, error)
	CheckGroupExists(ctx context.Context, userID, groupID uuid.UUID) (bool, error)
	GetStat(ctx context.Context, groupID uuid.UUID) (*models.GroupStat, error)
}
type MentorRepository interface {
	GetMyHelpers(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	UpdateRequest(ctx context.Context, request *models.HelpRequest) error
	GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
	CheckIsMentor(ctx context.Context, userID, groupID uuid.UUID) (bool, error)
	CheckRequest(ctx context.Context, id, mentorID uuid.UUID) (bool, error)
	GetGroup(ctx context.Context, userID, groupID uuid.UUID) (*models.Group, error)
}
type UsersRepository interface {
	Login(ctx context.Context, person *models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	GetByName(ctx context.Context, name string) (person *models.User, err error)
	EditUser(ctx context.Context, userID uuid.UUID, updates map[string]any) (*models.User, error)
	GetMyMentors(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
	GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.Role, error)
	GetMyRequests(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	CreateRequest(ctx context.Context, request *models.HelpRequest) error
	CheckExists(ctx context.Context, userID uuid.UUID) (bool, error)
	CheckIsStudent(ctx context.Context, userID, groupID uuid.UUID) (bool, error)
	CheckGroupExists(ctx context.Context, groupID uuid.UUID) (bool, error)
	GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error)
	AddRole(ctx context.Context, role *models.Role) error
	GetGroups(ctx context.Context, userID uuid.UUID, role string) ([]*models.Group, error)
}

type MinioRepository interface {
	UploadImage(file *models.File) (string, error)
	GetImage(image string) (string, error)
	DeleteImage(image string) error
}
