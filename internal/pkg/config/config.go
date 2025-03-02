package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	HttpPort             int    `envconfig:"HTTP_PORT" default:"8080"`
	ServerAddress        string `envconfig:"SERVER_ADDRESS" default:"0.0.0.0:8080"`
	ServerPort           int    `envconfig:"SERVER_PORT" default:"8080"`
	PostgresURL          string `envconfig:"POSTGRES_CONN" default:"postgresql://username:password@192.168.1.103:5432/database"`
	RedisHost            string `envconfig:"REDIS_HOST" default:"192.168.1.103"`
	RedisPort            int32  `envconfig:"REDIS_PORT" default:"6379"`
	RandomSecret         string `envconfig:"RANDOM_SECRET" default:"111"`
	BaseURL              string `envconfig:"BASE_URL" default:"http://localhost:8080"`
	TelegramBotToken     string `envconfig:"TELEGRAM_BOT_TOKEN" default:""`
	TelegramChatID       int64  `envconfig:"TELEGRAM_BOT_CHAT_ID" default:""`
	TelegramAdmins       []int64
	TelegramStringAdmins string `envconfig:"TELEGRAM_BOT_ADMINS" default:""`
	MinioHost            string `envconfig:"MINIO_HOST" default:"minio"`
	MinioPort            int    `envconfig:"MINIO_PORT" default:"9000"`
	MinioUser            string `envconfig:"MINIO_ROOT_USER" default:"cooluser"`
	MinioPassword        string `envconfig:"MINIO_ROOT_PASSWORD" default:"minio123"`
	BucketName           string `envconfig:"MINIO_BUCKET_NAME" default:"imagebucket"`
	CryptoKey            string `envconfig:"CRYPTO_KEY" default:"12345678901234567890123456789012"`
	IsModerationEnabled  bool   `envconfig:"IS_MODERATION_ENABLED" default:"false"`
	AIToken              string `envconfig:"AI_TOKEN" default:""`
	AIModel              string `envconfig:"AI_MODEL" default:""`
}

func New() (Config, error) {
	cfg := Config{}

	wd, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	var envPath string
	if strings.HasPrefix(wd, "/app") {
		wd = "/app"
		envPath = filepath.Join(wd, ".env")
	} else {
		wd = filepath.Join(wd)
		envPath = filepath.Join(wd, "dev.env")
	}

	fmt.Printf("Loading config from %s\n", envPath)

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return cfg, err
	}

	_ = godotenv.Load(envPath)

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	cfg.TelegramAdmins = make([]int64, 0, len(cfg.TelegramAdmins))
	for _, s := range strings.Split(cfg.TelegramStringAdmins, ",") {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			cfg.TelegramAdmins = append(cfg.TelegramAdmins, id)
		}
	}

	return cfg, nil
}
