package main

import (
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/xLeSHka/mentorLinkSchool/internal/app"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"go.uber.org/fx"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	_, _ = prettyjson.Marshal(cfg)

	fx.New(
		fx.Supply(cfg),
		app.WSApp,
	).Run()
}
