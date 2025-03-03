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

func TestCreateGroup(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Name         Name
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Name:         Name{"Group1"},
			jwt:          profile1JWT,
			name:         "create user",
			expectedCode: http.StatusOK,
		},
		{
			Name:         Name{"Group2"},
			jwt:          "1",
			name:         "empty jwt",
			expectedCode: http.StatusUnauthorized,
		},
	}
	db.Create(&profile1)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.Name)
			url := "/api/groups/create"
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
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
				err := db.Model(&models.Group{}).First(&group, test.Name).Error
				assert.Nil(t, err)
				assert.Equal(t, test.Name.Name, group.Name)
			}
		})
	}
}
