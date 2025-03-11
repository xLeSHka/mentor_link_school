package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestUpdateInviteCode(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Group        models.Group
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Group:        group1,
			jwt:          profile3JWT,
			name:         "update group",
			expectedCode: http.StatusOK,
		},
		{
			Group:        group1,
			jwt:          profile2JWT,
			name:         "failed update group",
			expectedCode: http.StatusForbidden,
		},
	}

	db.Create(&profile3)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleOwner)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/groups/" + test.Group.ID.String() + "/inviteCode"
			req, _ := http.NewRequest(http.MethodPost, url, nil)
			req.Header.Set("Authorization", "Bearer "+test.jwt) // bytes.NewBuffer(jsonData)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var group models.Group
				err = db.Model(&models.Group{}).First(&group).Error
				assert.Nil(t, err)
				assert.NotNil(t, group.InviteCode)
			}
		})
	}
}
