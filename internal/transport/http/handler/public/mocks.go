package publicRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *Route) mocks(c *gin.Context) {
	r.DB.FirstOrCreate(&models.Group{
		ID:        uuid.MustParse("5ad3e7ac-38da-4b0b-9bde-aa5f2050ad35"),
		AvatarURL: nil,
		Name:      "main",
	})
	r.DB.Create(&models.User{
		ID:        uuid.MustParse("17b015fc-0398-453f-bc0a-31bcf02b3ec1"),
		Name:      "student",
		AvatarURL: nil,
		BIO:       nil,
		Telegram:  "@student",
	})

	c.JSON(200, gin.H{
		"message": "PROOOOOOOOOOOOOOOOOD",
	})
}
