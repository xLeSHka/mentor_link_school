package service

import (
	"context"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"

	"github.com/google/uuid"
)

type GroupService interface {
	Create(ctx context.Context, group *models.Group, userID uuid.UUID) (string, error)
	AddRole(ctx context.Context, role *models.Role) error
	RemoveRole(ctx context.Context, role *models.Role) error
	UpdateInviteCode(ctx context.Context, groupID uuid.UUID) (string, error)
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]*models.User, error)
	Edit(ctx context.Context, group *models.Group) (*models.Group, error)
	GetStat(ctx context.Context, groupID uuid.UUID) (*models.GroupStat, error)
	UploadImage(ctx context.Context, file *models.File, groupID uuid.UUID) (string, *httpError.HTTPError)
	GetRoles(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Role, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
}
type MentorService interface {
	GetMyHelps(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error)
	UpdateRequest(ctx context.Context, request *models.HelpRequest) error
	GetStudents(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Pair, error)
}

type StudentService interface {
	CreateRequest(ctx context.Context, request *models.HelpRequest) error
	GetMentors(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Role, error)
	GetMyHelps(ctx context.Context, userID, groupID uuid.UUID) ([]*models.HelpRequest, error)
	GetMyMentors(ctx context.Context, userID, groupID uuid.UUID) ([]*models.Pair, error)
	GetRequestByID(ctx context.Context, reqID, groupID uuid.UUID) (*models.HelpRequest, error)
	//GetRequest(ctx context.Context, UserID, MentorID, GroupID uuid.UUID) (models.HelpRequest, error)
}
type UsersService interface {
	Login(ctx context.Context, telegram, password string) (string, error)
	Register(ctx context.Context, user *models.User) (string, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	GetByTelegram(ctx context.Context, telegram string) (person *models.User, err error)
	UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError)
	Edit(ctx context.Context, userID uuid.UUID, user *models.User) (*models.User, error)
	GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.GroupWithRoles, error)
	GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error)
	Invite(ctx context.Context, inviteCode string, userID uuid.UUID) (bool, error)
}
