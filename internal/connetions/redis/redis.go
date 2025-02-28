package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"
	"go.uber.org/fx"
	"strconv"
)

func New(config config.Config, lc fx.Lifecycle) (*redis.Client, error) {
	//println("Redis", config.RedisHost+":", config.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisHost + ":" + strconv.Itoa(int(config.RedisPort)),
		DB:   0,
	})

	err := rdb.Ping(context.Background()).Err()

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return rdb.Close()
		},
	})

	return rdb, nil
}
