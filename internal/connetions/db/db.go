package db

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(config config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	}

	db, err := gorm.Open(postgres.Open(config.PostgresURL), gormConfig)
	if err != nil {
		fmt.Println(err)
		panic("Error db connecting")
	}

	//db.AutoMigrate(&models.Company{})
	//db.AutoMigrate(&models.Promocode{})
	//db.AutoMigrate(&models.PromocodeCategory{})
	//db.AutoMigrate(&models.PromocodeUnique{})
	//db.AutoMigrate(&models.User{})
	//db.AutoMigrate(&models.PromocodeActivation{})

	//lc.Append(fx.Hook{
	//	OnStop: func(ctx context.Context) error {
	//		return db.Dis
	//	},
	//})

	return db, nil
}

func MigrationUp(cfg config.Config) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	//println(fmt.Sprintf("%s?sslmode=disable", os.Getenv("POSTGRES_CONN")))

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		return err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
