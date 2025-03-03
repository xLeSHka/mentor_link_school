package tests

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestGroupUploadImage(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		GroupID      uuid.UUID
		jwt          string
		name         string
		expectedCode int
		Formfile     string
		Imagename    string
	}
	tests := []Test{
		{
			GroupID:      group1.ID,
			jwt:          profile3JWT,
			name:         "add image to group1",
			expectedCode: http.StatusOK,
			Imagename:    "go.jpg",
			Formfile:     "image",
		},
		{
			GroupID:      group1.ID,
			jwt:          profile3JWT,
			name:         "failed upload gif",
			expectedCode: http.StatusBadRequest,
			Imagename:    "go.gif",
			Formfile:     "image",
		},
		{
			GroupID:      group1.ID,
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusForbidden,
			Imagename:    "go.jpg",
			Formfile:     "image",
		},
	}

	db.Create(&group1)
	db.Create(&profile3)
	db.Create(&roleOwner)
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

			req, _ := http.NewRequest(http.MethodPost, "/api/groups/"+test.GroupID.String()+"/uploadAvatar", buff)
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
