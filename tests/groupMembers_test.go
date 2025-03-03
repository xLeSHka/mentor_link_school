package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestGetGroupMembers(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []respGetMember
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected: []respGetMember{
				{
					UserID:    profile2.ID,
					AvatarUrl: profile2.AvatarURL,
					Name:      profile2.Name,
					Role:      "mentor",
				},
				{
					UserID:    profile1.ID,
					AvatarUrl: profile1.AvatarURL,
					Name:      profile1.Name,
					Role:      "student",
				},
			},
			jwt:          profile3JWT,
			name:         "succes group members",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusForbidden,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&accepted)
	db.Create(&pair)
	db.Create(&profile3)
	db.Create(&roleOwner)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/groups/" + group1.ID.String() + "/members"
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
				var resp []respGetMember
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, _ := range test.Expected {

					assert.Equal(t, test.Expected[i].UserID, resp[i].UserID)
					assert.Equal(t, test.Expected[i].Name, resp[i].Name)
					assert.Equal(t, test.Expected[i].Role, resp[i].Role)
				}
			}
		})
	}
}
