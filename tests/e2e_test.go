package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
			expectedCode: http.StatusBadRequest,
		},
		{
			Req: reqUpdateStatus{
				ID:     pending.ID,
				Status: true,
			},
			jwt:          profile1JWT,
			name:         "not own request",
			expectedCode: http.StatusBadRequest,
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

func TestMentorGetStudents(t *testing.T) {
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
			Expected:     []models.User{profile1},
			jwt:          profile2JWT,
			name:         "succes mentor get students",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.User{},
			jwt:          profile1JWT,
			name:         "empty mentor students",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&accepted)
	db.Create(&pair)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/mentors/students"
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
				var resp []respGetMyStudent
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].StudentID)
					assert.Equal(t, profile1.Name, resp[i].Name)
				}
			}
		})
	}
}

func TestStudentGetMentors(t *testing.T) {
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
			name:         "succes mentor get students",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.User{},
			jwt:          profile2JWT,
			name:         "empty mentor students",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&accepted)
	db.Create(&pair)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/user/mentors"
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
				var resp []respGetMyMentor
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].MentorID)
					assert.Equal(t, profile2.Name, resp[i].Name)
				}
			}
		})
	}
}

func TestMentorUploadImage(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		UserID       uuid.UUID
		jwt          string
		name         string
		expectedCode int
		Formfile     string
		Imagename    string
	}
	tests := []Test{
		{
			jwt:          profile2JWT,
			name:         "add image to profile2",
			expectedCode: http.StatusOK,
			Imagename:    "go.jpg",
			Formfile:     "image",
		},
		{
			jwt:          profile2JWT,
			name:         "failed upload gif",
			expectedCode: http.StatusBadRequest,
			Imagename:    "go.gif",
			Formfile:     "image",
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
			Imagename:    "go.jpg",
			Formfile:     "image",
		},
	}
	db.Create(&profile2)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			image, err := os.Open("./image/" + test.Imagename)
			assert.Nil(t, err)
			defer image.Close()

			buff := &bytes.Buffer{}
			buffWriter := io.Writer(buff)

			ext := filepath.Ext(test.Imagename)
			imageURL := profile2.ID.String() + ext

			formWriter := multipart.NewWriter(buffWriter)

			formPart, err := formWriter.CreateFormFile(test.Formfile, imageURL)
			assert.Nil(t, err)
			_, err = io.Copy(formPart, image)
			assert.Nil(t, err)
			formWriter.Close()

			req, _ := http.NewRequest(http.MethodPost, "/api/user/uploadAvatar", buff)
			req.Header.Set("Content-Type", formWriter.FormDataContentType())
			req.Header.Set("Authorization", "Bearer "+test.jwt)

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				_, err := MinioRepository.GetImage(imageURL)
				assert.Nil(t, err)
				ext := filepath.Ext(test.Imagename)
				imagename := fmt.Sprintf("%s%s", profile2.ID.String(), ext)

				MinioRepository.DeleteImage(imagename)
			}
		})
	}
}
