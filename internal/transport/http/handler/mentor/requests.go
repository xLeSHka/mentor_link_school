package mentorsRoute

import "gitlab.prodcontest.ru/team-14/lotti/internal/models"

type GetGroupID struct {
	ID string `uri:"groupId" binding:"required,uuid"`
}
type GetMentorID struct {
	ID      string `uri:"mentorId" binding:"required,uuid"`
	GroupID string `uri:"groupId" binding:"required,uuid"`
}
type GetMentorRequestDto struct {
	Goal string `string:"goal" binding:"required"`
}

func mapMentor(mentor *models.Mentor) {

}
