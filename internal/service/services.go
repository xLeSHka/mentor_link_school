package service

import (
	"context"

	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

type GroupService interface {
	Create(ctx context.Context, group *models.Group, userID uuid.UUID) error
	UpdateRole(ctx context.Context, ownerID, groupID, userID uuid.UUID, role string) error
	UpdateInviteCode(ctx context.Context, groupID, ownerID uuid.UUID) (string, error)
	GetMembers(ctx context.Context, ownerID, groupID uuid.UUID) ([]*models.Role, error)
	GetStat(ctx context.Context, ownerID, groupID uuid.UUID) (*models.GroupStat, error)
	UploadImage(ctx context.Context, file *models.File, groupID, personID uuid.UUID) (string, *httpError.HTTPError)
}
type MentorService interface {
	GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequestWithGIDs, error)
	UpdateRequest(ctx context.Context, request *models.HelpRequest) error
	GetStudents(ctx context.Context, userID uuid.UUID) ([]*models.PairWithGIDs, error)
}

type UserService interface {
	Login(ctx context.Context, name string) (*models.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error)
	UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError)
	GetMyMentors(ctx context.Context, userID uuid.UUID) ([]*models.PairWithGIDs, error)
	GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequestWithGIDs, error)
	CreateRequest(ctx context.Context, request *models.HelpRequest) error
	GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.RoleWithGIDs, error)
	Invite(ctx context.Context, inviteCode string, userID uuid.UUID) (bool, error)
	GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Role, error)
	GetCommonGroups(userID, mentorID uuid.UUID) ([]uuid.UUID, error)
	GetRequestByID(ctx context.Context, reqID uuid.UUID) (models.HelpRequest, error)
	GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error)
	GetGroupByID(ctx context.Context, ID uuid.UUID) (*models.Group, error)
}
