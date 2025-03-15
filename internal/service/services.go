package service

import (
	"context"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"

	"github.com/google/uuid"
)

type GroupService interface {
	Create(ctx context.Context, group *models.Group, userID uuid.UUID) error
	AddRole(ctx context.Context, role *models.Role) error
	RemoveRole(ctx context.Context, role *models.Role) error
	UpdateInviteCode(ctx context.Context, groupID, personID uuid.UUID) (string, error)
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]*models.User, error)
	GetStat(ctx context.Context, groupID uuid.UUID) (*models.GroupStat, error)
	Edit(ctx context.Context, userID, groupID uuid.UUID, updates map[string]any) (*models.Group, error)
	UploadImage(ctx context.Context, file *models.File, groupID, personID uuid.UUID) (string, *httpError.HTTPError)
	GetRoles(ctx context.Context, userID uuid.UUID) ([]*models.Role, error)
	BanUser(ctx context.Context, userID uuid.UUID) error
	GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error)
	GetGroupByID(ctx context.Context, ID uuid.UUID) (*models.Group, error)
	GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Role, error)
	Invite(ctx context.Context, inviteCode string, userID uuid.UUID) (bool, error)
}
type MentorService interface {
	GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	UpdateRequest(ctx context.Context, request *models.HelpRequest) error
	GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
}

type UserService interface {
	Login(ctx context.Context, telegram, password string) (string, error)
	Register(ctx context.Context, user *models.User) (string, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError)
	GetMyMentors(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error)
	GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error)
	CreateRequest(ctx context.Context, request *models.HelpRequest) error
	GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.Role, error)
	GetRequestByID(ctx context.Context, reqID uuid.UUID) (*models.HelpRequest, error)
	Edit(ctx context.Context, userID uuid.UUID, user *models.User) (*models.User, error)
}
