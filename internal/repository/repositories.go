package repository

import (
	"context"

	"github.com/xLeSHka/mentorLinkSchool/internal/models"

	"github.com/google/uuid"
)

type CacheRepository interface {
	AddRoles(ctx context.Context, roles []*models.Role) error
	RemoveRoles(ctx context.Context, roles []*models.Role) error
	SaveToken(ctx context.Context, userID uuid.UUID, token string) error
	DeleteToken(ctx context.Context, userID uuid.UUID) error
}
type GroupRepository interface {
	AddRole(ctx context.Context, role *models.Role) error
	RemoveRole(ctx context.Context, role *models.Role) error
	Create(ctx context.Context, group *models.Group, userID uuid.UUID) error
	UpdateInviteCode(ctx context.Context, groupID uuid.UUID, inviteCode string) error
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]*models.User, error)
	//CheckGroupExists(ctx context.Context, userID, groupID uuid.UUID) (bool, error)
	Edit(ctx context.Context, group *models.Group) (*models.Group, error)
	GetGroupByID(ctx context.Context, ID uuid.UUID) (*models.Group, error)
	GetStat(ctx context.Context, groupID uuid.UUID) (*models.GroupStat, error)

	GetRoles(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Role, error)
}
type MentorRepository interface {
	UpdateRequest(ctx context.Context, request *models.HelpRequest) error
	GetStudents(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Pair, error)
	GetMyHelpers(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error)
	CreatePair(ctx context.Context, pair *models.Pair) error
	CheckIsMentor(ctx context.Context, userID, groupID uuid.UUID) (bool, error)
	//CheckRequest(ctx context.Context, id, mentorID uuid.UUID) (bool, error)
	//GetRequest(ctx context.Context, UserID, MentorID, GroupID uuid.UUID) (models.HelpRequest, error)
}
type UsersRepository interface {
	Login(ctx context.Context, telegram string) (*models.User, error)
	Register(ctx context.Context, person *models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	EditUser(ctx context.Context, userID uuid.UUID, user *models.User) (*models.User, error)
	GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Roles, error)
	GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error)
}

type MinioRepository interface {
	UploadImage(file *models.File) (string, error)
	GetImage(image string) (string, error)
	DeleteImage(image string) error
}
type StudentRepository interface {
	GetMyMentors(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Pair, error)
	GetMentors(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Role, error)
	GetMyRequests(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error)
	CreateRequest(ctx context.Context, request *models.HelpRequest) error
	GetRequestByID(ctx context.Context, reqID, groupID uuid.UUID) (*models.HelpRequest, error)
}
