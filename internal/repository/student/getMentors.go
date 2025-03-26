package repositoryStudent

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
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

func (r *StudentRepository) GetMentors(ctx context.Context, userID, groupID uuid.UUID, page, size int) ([]*models.Role, int64, error) {
	var role []*models.Role
	err := r.DB.WithContext(ctx).Model(&models.Role{}).
		Where("role = 'mentor' AND group_id = ? AND user_id != ?", groupID, userID).
		Where("NOT EXISTS (SELECT 1 FROM help_requests WHERE user_id = ? AND group_id = ? AND mentor_id = roles.user_id AND (status = 'pending' OR status = 'accepted'))", userID, groupID).
		Preload("User").
		Offset(page * size).Limit(size).
		Find(&role).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = r.DB.Model(&models.Role{}).Where("role = 'mentor' AND group_id = ? AND user_id != ?", groupID, userID).
		Where("NOT EXISTS (SELECT 1 FROM help_requests WHERE user_id = ? AND group_id = ? AND mentor_id = roles.user_id AND (status = 'pending' OR status = 'accepted'))", userID, groupID).
		Count(&count).Error
	return role, count, err
}
