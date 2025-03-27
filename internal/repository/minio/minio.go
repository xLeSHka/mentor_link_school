package repositoryMinio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"strings"
	"time"
)

func (r *MinioRepository) UploadImage(file *models.File) (string, error) {
	err := r.MC.RemoveObject(context.Background(), r.BN, file.Filename, minio.RemoveObjectOptions{})
	if err != nil {
		return "", err
	}
	_, err = r.MC.PutObject(context.Background(), r.BN, file.Filename, file.File, file.Size, minio.PutObjectOptions{ContentType: file.Mimetype})
	if err != nil {
		return "", err
	}
	url, err := r.MC.PresignedGetObject(context.Background(), r.BN, file.Filename, time.Hour*24*7, nil)
	if err != nil {
		return "", err
	}
	avatarURL := url.String()
	avatarURL = strings.Replace(avatarURL, "http://minio:9000", "https://localhost", 1)
	avatarURL = strings.Split(avatarURL, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	return avatarURL, nil
}
func (r *MinioRepository) GetImage(image string) (string, error) {
	url, err := r.MC.PresignedGetObject(context.Background(), r.BN, image, time.Hour*24*7, nil)
	if err != nil {
		return "", err
	}
	avatarURL := url.String()
	avatarURL = strings.Replace(avatarURL, "http://minio:9000", "https://localhost", 1)
	avatarURL = strings.Split(avatarURL, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	return avatarURL, nil
}
func (r *MinioRepository) DeleteImage(filename string) error {
	err := r.MC.RemoveObject(context.Background(), r.BN, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
