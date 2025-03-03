package repositoryUser

import (
	"github.com/google/uuid"
)

func (r *UsersRepository) GetCommonGroups(userID, mentorID uuid.UUID) ([]uuid.UUID, error) {
	var groupIDs []uuid.UUID
	err := r.DB.Table("roles").Select("group_id").Where("user_id = ?", userID).Find(&groupIDs).Error
	var intersection []uuid.UUID
	err = r.DB.Table("roles").Select("group_id").Where("role = 'mentor' OR role = 'student-mentor' AND group_id IN (?)", groupIDs).Find(&intersection).Error
	return intersection, err
}
