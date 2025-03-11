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

func TestEditUser(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Req          reqEditUser
		jwt          string
		name         string
		expectedCode int
	}
	var r = reqEditUser{
		Name:     "New profile 1",
		BIO:      "New bio",
		Telegram: "new telegram",
	}
	tests := []Test{
		{
			Req:          r,
			jwt:          profile1JWT,
			name:         "get profile 1",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "bad request",
			expectedCode: http.StatusBadRequest,
		},
	}
	db.Create(&profile1)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.Req)
			assert.Nil(t, err)
			url := "/api/user/profile/edit"
			req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData)) // bytes.NewBuffer(jsonData)
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
				var resp models.User
				err = db.Model(&models.User{}).Where("id = ?", profile1.ID).First(&resp).Error
				assert.Nil(t, err)
				assert.Equal(t, r.Name, resp.Name)
				assert.Equal(t, r.BIO, *resp.BIO)
				assert.Equal(t, r.Telegram, resp.Telegram)
			}
		})
	}
}
