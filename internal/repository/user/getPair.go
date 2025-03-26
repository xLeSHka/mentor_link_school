package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) GetPair(ctx context.Context, userID, userID2, groupID uuid.UUID) (*models.Pair, error) {
	var pair models.Pair
	err := r.DB.Model(&models.Pair{}).Where("user_id IN (?,?) AND mentor_id IN (?,?) AND group_id = ?", userID, userID2, userID, userID2, groupID).First(&pair).Error
	return &pair, err
}
