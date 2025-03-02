package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestLogin(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Name         Name
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Name:         Name{"User1"},
			name:         "create user",
			expectedCode: http.StatusOK,
		},
		{
			Name:         Name{"User1"},
			name:         "login user",
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.Name)
			url := "/api/user/auth/sign-in"
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData)) // bytes.NewBuffer(jsonData)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var user models.User
				err := db.Model(&models.User{}).First(&user, test.Name).Error
				assert.Nil(t, err)
				assert.Equal(t, test.Name.Name, user.Name)
			}
		})
	}
}
