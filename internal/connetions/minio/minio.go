package minio

import (
	"context"
	"fmt"
	"gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New(config config.Config) (*minio.Client, error) {
	address := fmt.Sprintf("%s:%d", config.MinioHost, config.MinioPort)

	client, err := minio.New(address, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioUser, config.MinioPassword, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		policy := fmt.Sprintf("{\n    \"Version\": \"2012-10-17\",\n    \"Statement\": [\n              {\n            \"Effect\": \"Allow\",\n            \"Principal\": {\n                \"AWS\": [\n                    \"*\"\n                ]\n            },\n            \"Action\": [\n                \"s3:GetObject\"\n            ],\n            \"Resource\": [\n                \"arn:aws:s3:::%s/*\"\n            ]\n        }\n    ]\n}", config.BucketName)
		err := client.SetBucketPolicy(context.Background(), config.BucketName, policy)
		if err != nil {
			log.Println("Error setting bucket policy", err.Error())
		}
		log.Println("Successfully set bucket policy")
	}()
	exists, err := client.BucketExists(context.Background(), config.BucketName)
	if err != nil {
		return nil, err
	}
	if exists {
		return client, nil
	}
	err = client.MakeBucket(context.Background(), config.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		return nil, err
	}
	return client, nil
}
