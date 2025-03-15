package cache

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"time"
)

type RedisRepository struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}
func (c *RedisRepository) AddRole(ctx context.Context, roles []*models.Role) error {
	for _, role := range roles {
		err := c.rdb.SAdd(ctx, "roles:"+uuid.New().String(), role).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *RedisRepository) RemoveRoles(ctx context.Context, roles []*models.Role) error {
	for _, role := range roles {
		err := c.rdb.SRem(ctx, "roles:"+uuid.New().String(), role).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *RedisRepository) SaveToken(ctx context.Context, userID uuid.UUID, token string) error {
	err := c.rdb.Set(ctx, "jwt:"+userID.String(), token, time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c *RedisRepository) DeleteToken(ctx context.Context, userID uuid.UUID) error {
	err := c.rdb.Del(ctx, "jwt:"+userID.String()).Err()
	if err != nil {
		return err
	}
	return nil
}
