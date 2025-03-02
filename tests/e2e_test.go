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
					assert.Equal(t, group1.ID, user[i].GroupID)
				}
			}
		})
	}
}
func TestCreateRequest(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Request      reqCreateHelp
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Request: reqCreateHelp{
				Requests: []Pair{
					{
						MentorID: profile2.ID,
						GroupId:  group1.ID,
					},
				},
				Goal: "Some goal",
			},
			jwt:          profile1JWT,
			name:         "succes create",
			expectedCode: http.StatusOK,
		},
		{
			Request: reqCreateHelp{
				Requests: []Pair{
					{
						MentorID: uuid.New(),
						GroupId:  group1.ID,
					},
				},
				Goal: "Some goal",
			},
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Request: reqCreateHelp{
				Requests: []Pair{
					{
						MentorID: profile1.ID,
						GroupId:  group1.ID,
					},
				},
				Goal: "Some goal",
			},
			jwt:          profile2JWT,
			name:         "user is not student",
			expectedCode: http.StatusBadRequest,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/group/requests"
			jsonData, _ := json.Marshal(&test.Request)
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
		})
	}
}
func TestUserGetRequests(t *testing.T) {
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
			jwt:          profile1JWT,
			name:         "succes student get requests",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.HelpRequest{},
			jwt:          profile2JWT,
			name:         "empty student requests",
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
			url := "/api/user/requests"
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
				var resp []respGetHelp
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].ID)
					assert.Equal(t, exp.MentorID, resp[i].MentorID)
					assert.Equal(t, profile2.Name, resp[i].MentorName)
					assert.Equal(t, exp.Goal, resp[i].Goal)
					assert.Equal(t, exp.Status, resp[i].Status)
				}
			}
		})
	}
}

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
