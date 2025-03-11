package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestMentorUpdateRequests(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     models.HelpRequest
		Req          reqUpdateStatus
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected: accepted,
			Req: reqUpdateStatus{
				ID:     pending.ID,
				Status: true,
			},
			jwt:          profile2JWT,
			name:         "succes mentor update request",
			expectedCode: http.StatusOK,
		},
		{
			Req: reqUpdateStatus{
				ID:     pending.ID,
				Status: true,
			},
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Req: reqUpdateStatus{
				ID:     pending.ID,
				Status: true,
			},
			jwt:          profile1JWT,
			name:         "not own request",
			expectedCode: http.StatusNotFound,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&pending)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/mentors/requests"
			jsonData, _ := json.Marshal(test.Req)
			assert.Nil(t, err)
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData)) // bytes.NewBuffer(jsonData)
			req.Header.Set("Authorization", "Bearer "+test.jwt)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var resp models.HelpRequest
				err = db.Model(&models.HelpRequest{}).Where("id = ?", test.Req.ID).First(&resp).Error
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.Status, resp.Status)
				var p models.Pair
				err = db.Model(&models.Pair{}).Where("user_id = ? AND mentor_id = ? AND group_id = ?", resp.UserID, resp.MentorID, resp.GroupID).First(&p).Error
				assert.Nil(t, err)
			}
		})
	}
}
