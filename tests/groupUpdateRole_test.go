package tests

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestUpdateRole(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		GroupID      uuid.UUID
		Req          reqUpdateRole
		jwt          string
		name         string
		expectedCode int
	}
	r := reqUpdateRole{
		Role: "mentor",
		ID:   profile1.ID.String(),
	}
	tests := []Test{
		{
			GroupID:      group1.ID,
			jwt:          profile3JWT,
			name:         "update role profile 1",
			expectedCode: http.StatusOK,
			Req:          r,
		},
		{
			GroupID:      group1.ID,
			jwt:          profile2JWT,
			name:         "failed update role profile 1",
			expectedCode: http.StatusForbidden,
			Req:          r,
		},
		{
			GroupID:      group1.ID,
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusForbidden,
			Req:          r,
		},
	}
	db.Create(&profile1)
	db.Create(&profile3)
	db.Create(&group1)
	db.Create(&roleOwner)
	db.Create(&roleStudent)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.Req)
			assert.Nil(t, err)
			req, _ := http.NewRequest(http.MethodPost, "/api/groups/"+test.GroupID.String()+"/members/role", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
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
				err = db.Model(&models.Role{}).Where("user_id = ? AND group_id = ?", test.Req.ID, test.GroupID).First(&role).Error
				assert.Nil(t, err)
				assert.Equal(t, test.Req.Role, role.Role)
			}
		})
	}
}
