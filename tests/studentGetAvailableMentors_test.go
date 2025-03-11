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

func TestAvailableMentors(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []models.User
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     []models.User{profile2},
			jwt:          profile1JWT,
			name:         "get available profile 2",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "get unknown user available",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.User{},
			jwt:          profile2JWT,
			name:         "available mentors not found",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/user/availableMentors"
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
				var user []*respGetMentor
				err = json.Unmarshal(w.Body.Bytes(), &user)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.Name, user[i].Name)
					assert.Equal(t, *exp.BIO, *user[i].BIO)
					assert.Equal(t, exp.ID, user[i].MentorID)
				}
			}
		})
	}
}
