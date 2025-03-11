package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestMentorGetRequests(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []models.HelpRequest
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     []models.HelpRequest{pending},
			jwt:          profile2JWT,
			name:         "succes mentor get requests",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.HelpRequest{},
			jwt:          profile1JWT,
			name:         "empty mentor requests",
			expectedCode: http.StatusOK,
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
			assert.Nil(t, err)
			req, _ := http.NewRequest(http.MethodGet, url, nil) // bytes.NewBuffer(jsonData)
			req.Header.Set("Authorization", "Bearer "+test.jwt)

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var resp []respGetRequest
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].ID)
					assert.Equal(t, exp.UserID, resp[i].UserID)
					assert.Equal(t, profile1.Name, resp[i].Name)
					assert.Equal(t, exp.Goal, resp[i].Goal)
					assert.Equal(t, exp.Status, resp[i].Status)
				}
			}
		})
	}
}
