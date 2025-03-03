package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
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
		Expected     models.GroupStat
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     stats,
			jwt:          profile3JWT,
			name:         "success get stat",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
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
			url := "/api/groups/" + group1.ID.String() + "/stat"
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
				var resp respStat
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.StudentsCount, resp.StudentsCount)
				assert.Equal(t, test.Expected.MentorsCount, resp.MentorsCount)
				assert.Equal(t, test.Expected.HelpRequestCount, resp.HelpRequestCount)
				assert.Equal(t, test.Expected.AcceptedRequestCount, resp.AcceptedRequestCount)
				assert.Equal(t, test.Expected.RejectedRequestCount, resp.RejectedRequestCount)
				assert.Equal(t, test.Expected.Conversion, resp.Conversion)
			}
		})
	}
}
