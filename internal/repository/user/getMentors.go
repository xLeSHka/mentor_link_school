package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

//func (r *UsersRepository) GetAvaliableMentors(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error) {
//	var resp []*models.Pair
//	err := r.DB.
//		WithContext(ctx).
//		Where("group_id in (select group_id from roles where user_id = ? AND role = 'student')", userID).
//		Preload("User").
//		Table("users").Find(&resp).Error
//
//	return resp, err
//
//}

func (r *UsersRepository) GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.RoleWithGIDs, error) {
	var role []*models.RoleWithGIDs
	err := r.DB.WithContext(ctx).Table("roles").
		Select("user_id,array_agg(group_id) as group_ids").
		Where("role = 'mentor' OR role = 'student-mentor' AND group_id in (SELECT group_id FROM roles WHERE user_id = ? AND role = 'student' OR role = 'student-mentor')", userID).
		Where("group_id in (SELECT group_id FROM roles WHERE user_id = ?)", userID).
		Group("user_id").
		Where("NOT EXISTS (SELECT 1 FROM pairs WHERE user_id = ? and mentor_id = roles.user_id)", userID).
		Where("NOT EXISTS (SELECT 1 FROM help_requests WHERE user_id = ? and mentor_id = roles.user_id AND status = 'pending' OR status = 'accepted')").
		Preload("User").
		Find(&role).Error
	if err != nil {
		return nil, err
	}
	return role, err
}
