package repositoryMinio

import (
	"github.com/minio/minio-go/v7"
	"gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
)

type MinioRepository struct {
	MC *minio.Client
	BN string
}

func New(mc *minio.Client, config config.Config) repository.MinioRepository {
	return &MinioRepository{
		mc,
		config.BucketName,
	}
}
