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

func TestGetProfile(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     models.User
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     profile1,
			jwt:          profile1JWT,
			name:         "get profile 1",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "get unknown profile",
			expectedCode: http.StatusNotFound,
		},
	}
	db.Create(&profile1)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/user/profile"
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
				var user resGetProfile
				err := json.Unmarshal(w.Body.Bytes(), &user)
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.Name, user.Name)
				assert.Equal(t, *test.Expected.BIO, *user.BIO)
			}
		})
	}
}
