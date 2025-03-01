package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

//func (r *UsersRepository) GetAvaliableMentors(ctx context.Context, userID uuid.UUID) {
//	var resp []*models.User
//	r.DB.
//}

func (r *UsersRepository) GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	return make([]*models.User, 0), nil
}
