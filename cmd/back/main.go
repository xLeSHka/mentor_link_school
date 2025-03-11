package main

import (
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/xLeSHka/mentorLinkSchool/internal/app"
	"github.com/xLeSHka/mentorLinkSchool/internal/connetions/db"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"go.uber.org/fx"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	s, _ := prettyjson.Marshal(cfg)
	fmt.Println(string(s))

	err = db.MigrationUp(cfg)
	if err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		return
	}
	fx.New(
		fx.Supply(cfg),
		app.App,
	).Run()
}
