package repositoryMinio

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"prodapp/internal/models"
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
	return url.String(), nil
}
func (r *MinioRepository) GetImage(image string) (string, error) {
	url, err := r.MC.PresignedGetObject(context.Background(), r.BN, image, time.Hour*24*7, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
func (r *MinioRepository) DeleteImage(personID uuid.UUID) error {
	err := r.MC.RemoveObject(context.Background(), r.BN, personID.String(), minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
