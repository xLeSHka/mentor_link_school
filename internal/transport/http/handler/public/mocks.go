package publicRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *Route) mocks(c *gin.Context) {
	owner := &models.User{
		ID:        uuid.MustParse("18015fc-0398-453f-bc0a-31bcf02b3ec1"),
		Name:      "owner 1",
		AvatarURL: nil,
		BIO:       nil,
		Telegram:  "@owner",
	}
	ownerGroup := &models.Group{
		ID:        uuid.MustParse("5ad3e7ac-38da-4b0b-9bde-aa5f2050ad35"),
		AvatarURL: nil,
		Name:      "owner group 1",
	}
	student := &models.User{
		ID:        uuid.MustParse("17b015fc-0398-453f-bc0a-31bcf02b3ec1"),
		Name:      "student 1",
		AvatarURL: nil,
		BIO:       nil,
		Telegram:  "@student",
	}
	group1 := &models.Group{
		ID:        uuid.MustParse("17b015fc-0398-453f-bc0a-31bcf02b3ec2"),
		AvatarURL: nil,
		Name:      "group 1",
	}
	group2 := &models.Group{
		ID:        uuid.MustParse("6ad3e7ac-38da-4b0b-9bde-aa5f2050ad35"),
		AvatarURL: nil,
		Name:      "group 2",
	}
	bio := "new bio"
	mentor1 := &models.User{
		ID:        uuid.MustParse("18015fc-0398-453f-bc0a-31bcf02b3ec1"),
		Name:      "mentor 1",
		AvatarURL: nil,
		BIO:       &bio,
		Telegram:  "@mentor",
	}
	helpReq1 := &models.HelpRequest{
		ID:       uuid.MustParse("20015fc-4398-453f-bc0a-31bcf02b3ec1"),
		UserID:   student.ID,
		MentorID: mentor1.ID,
		GroupID:  group2.ID,
		Goal:     "PROD Project",
		BIO:      student.AvatarURL,
		Status:   "pending",
	}
	helpReq2 := &models.HelpRequest{
		ID:       uuid.MustParse("30015fc-4398-453f-bc0a-31bcf02b3ec1"),
		UserID:   student.ID,
		MentorID: mentor1.ID,
		GroupID:  group1.ID,
		Goal:     "PRODANO Project",
		BIO:      student.AvatarURL,
		Status:   "accepted",
	}
	r.DB.FirstOrCreate(group1)
	r.DB.FirstOrCreate(student)
	r.DB.FirstOrCreate(&models.Role{
		UserID:  student.ID,
		GroupID: group1.ID,
		Role:    "student",
	})
	r.DB.FirstOrCreate(mentor1)
	r.DB.FirstOrCreate(&models.Role{
		UserID:  mentor1.ID,
		GroupID: group1.ID,
		Role:    "mentor",
	})
	r.DB.FirstOrCreate(&models.Role{
		UserID:  student.ID,
		GroupID: group2.ID,
		Role:    "student",
	})
	r.DB.FirstOrCreate(&models.Role{
		UserID:  mentor1.ID,
		GroupID: group2.ID,
		Role:    "mentor",
	})
	r.DB.FirstOrCreate(helpReq1)
	r.DB.FirstOrCreate(helpReq2)
	r.DB.FirstOrCreate(&models.Pair{
		UserID:   student.ID,
		MentorID: mentor1.ID,
		GroupID:  group1.ID,
		Goal:     "PRODANO Project",
	})

	r.DB.FirstOrCreate(owner)
	r.DB.FirstOrCreate(ownerGroup)
	r.DB.FirstOrCreate(&models.Role{
		UserID:  owner.ID,
		GroupID: ownerGroup.ID,
		Role:    "owner",
	})
	c.JSON(200, gin.H{
		"message": "PROOOOOOOOOOOOOOOOOD",
	})
}
