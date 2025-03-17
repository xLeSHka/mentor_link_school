package tg

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/tg/botkit"
	"go.uber.org/fx"
	"log"
)

func Run(bot *botkit.Bot, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := bot.Run()
				if err != nil {
					log.Println("failed to start bot", err)
				}
			}()
			log.Println("started bot")
			return nil
		},
		OnStop: nil,
	})
}
