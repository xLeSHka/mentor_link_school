package conn

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config2 "github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
)

func New(config config2.Config) (*tgbotapi.BotAPI, error) {
	botAPI, err := tgbotapi.NewBotAPI(config.BOT_API_TOKEN)
	if err != nil {
		return nil, err
	}
	return botAPI, nil
}
