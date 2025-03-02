package group

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"github.com/google/uuid"
)

func (r *GroupRepository) GetStat(ctx context.Context, groupID uuid.UUID) (*models.GroupStat, error) {

	var stat models.GroupStat
	err := r.DB.WithContext(ctx).Table("roles").Where("group_id = ?", groupID).Count(&stat.TotalCount).Error
	if err != nil {
		return nil, err
	}
	err = r.DB.WithContext(ctx).Table("roles").Where("group_id = ? AND role = 'student'", groupID).Count(&stat.StudentsCount).Error
	if err != nil {
		return nil, err
	}
	err = r.DB.WithContext(ctx).Table("roles").Where("group_id = ? AND role = 'mentor'", groupID).Count(&stat.MentorsCount).Error
	if err != nil {
		return nil, err
	}
	err = r.DB.Model(&models.HelpRequest{}).Where("group_id = ?", groupID).Count(&stat.HelpRequestCount).Error
	if err != nil {
		return nil, err
	}
	err = r.DB.Model(&models.HelpRequest{}).Where("group_id = ? AND status = 'accepted' ", groupID).Count(&stat.AcceptedRequestCount).Error
	if err != nil {
		return nil, err
	}
	err = r.DB.Model(&models.HelpRequest{}).Where("group_id = ? AND status = 'rejected", groupID).Count(&stat.RejectedRequestCount).Error
	if err != nil {
		return nil, err
	}
	if stat.HelpRequestCount == 0 {
		stat.Conversion = 0
	} else {
		stat.Conversion = float64(stat.AcceptedRequestCount) / float64(stat.HelpRequestCount) * 100
	}
	return &stat, nil
}
