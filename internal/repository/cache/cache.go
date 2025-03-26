package cache

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"strconv"
	"time"
)

type RedisRepository struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) repository.CacheRepository {
	return &RedisRepository{rdb: rdb}
}
func (c *RedisRepository) AddRoles(ctx context.Context, roles []*models.Role) error {
	for _, role := range roles {
		err := c.rdb.SAdd(ctx, "roles:"+role.UserID.String()+"_"+role.GroupID.String(), role.Role).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *RedisRepository) RemoveRoles(ctx context.Context, roles []*models.Role) error {
	for _, role := range roles {
		err := c.rdb.SRem(ctx, "roles:"+role.UserID.String()+"_"+role.GroupID.String(), role.Role).Err()
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
func (c *RedisRepository) SaveID(ctx context.Context, userID uuid.UUID, id int64) error {
	err := c.rdb.Set(ctx, "id:"+userID.String(), strconv.FormatInt(id, 10), time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c *RedisRepository) GetID(ctx context.Context, userID uuid.UUID) (int64, error) {
	res := c.rdb.Get(ctx, "id:"+userID.String())
	if res.Err() != nil {
		return 0, res.Err()
	}
	return strconv.ParseInt(res.Val(), 10, 64)
}
