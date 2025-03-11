package repositoryMinio

import (
	"github.com/minio/minio-go/v7"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
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
