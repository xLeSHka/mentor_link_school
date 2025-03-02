package group

import (
	"context"

	"github.com/google/uuid"
)

func (r *GroupRepository) GetGroupUsersStat(ctx context.Context, groupID uuid.UUID) error {
	type GroupStat struct {
		StudentsCount int64
		MentorsCount  int64
		TotalCount    int64
	}

	var stat GroupStat
	err := r.DB.WithContext(ctx).Table("roles").Where("group_id = ?", groupID).Count(&stat.TotalCount).Error
	if err != nil {
		return err
	}
	err = r.DB.WithContext(ctx).Table("roles").Where("group_id = ? AND role = 'student'", groupID).Count(&stat.StudentsCount).Error
	if err != nil {
		return err
	}
	err = r.DB.WithContext(ctx).Table("roles").Where("group_id = ? AND role = 'mentor'", groupID).Count(&stat.MentorsCount).Error
	if err != nil {
		return err
	}
	return nil
}
