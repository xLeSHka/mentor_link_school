package tests

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestEditGroup(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		GroupID      uuid.UUID
		Req          reqEditGroup
		jwt          string
		name         string
		expectedCode int
	}
	var r = reqEditGroup{
		Name: "New name 1",
	}
	tests := []Test{
		{
			GroupID:      group1.ID,
			Req:          r,
			jwt:          profile3JWT,
			name:         "edit group 1",
			expectedCode: http.StatusOK,
		},
		{
			GroupID:      group1.ID,
			Req:          reqEditGroup{},
			jwt:          unknownJWT,
			name:         "bad request",
			expectedCode: http.StatusBadRequest,
		},
	}
	db.Create(&group1)
	db.Create(&group2)
	db.Create(&profile3)
	db.Create(&roleOwner)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.Req)
			assert.Nil(t, err)
			url := "/api/groups/" + test.GroupID.String() + "/edit"
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
				var resp models.Group
				err = db.Model(&models.Group{}).Where("id = ?", test.GroupID).First(&resp).Error
				assert.Nil(t, err)
				assert.Equal(t, r.Name, resp.Name)
			}
		})
	}
}
