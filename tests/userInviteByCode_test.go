package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestInviteByCode(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		InviteCode   string
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			InviteCode:   *group2.InviteCode,
			jwt:          profile2JWT,
			name:         "profile 2 join to group2",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&group2)
	db.Create(&profile2)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/api/groups/join/"+test.InviteCode, nil)
			req.Header.Set("Authorization", "Bearer "+test.jwt)

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var role models.Role
				err = db.Model(&models.Role{}).Where("user_id = ? AND group_id = ?", profile2.ID, group2.ID).First(&role).Error
				assert.Nil(t, err)
			}
		})
	}
}
