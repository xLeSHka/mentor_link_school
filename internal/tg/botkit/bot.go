package botkit

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"github.com/xLeSHka/mentorLinkSchool/internal/service"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"log"
)

type Run func(CallStack) CallStack
type CallStack struct {
	ChatID  int64
	Bot     *Bot
	Update  *tgbotapi.Update
	Action  Run
	IsPrint bool
	Parent  *CallStack
	Data    string
}
type Data struct {
	User    *models.User
	Group   *models.Group
	Profile *models.User
	Size    int
	Page    int
	LastMes int
}

// Данные на время жизни приложения
var userRuns = map[int64]CallStack{}
var userDatas = map[int64]*Data{}

type Bot struct {
	DB              *gorm.DB
	CryptoKey       []byte
	Api             *tgbotapi.BotAPI
	MinioRepository repository.MinioRepository
	UsersService    service.UsersService
	StudentService  service.StudentService
	MentorService   service.MentorService
	GroupService    service.GroupService
	CacheRepository repository.CacheRepository
	UserRepository  repository.UsersRepository
}
type FxOpts struct {
	fx.In
	DB              *gorm.DB
	Api             *tgbotapi.BotAPI
	CacheRepository repository.CacheRepository
	MinioRepository repository.MinioRepository
	StudentService  service.StudentService
	UsersService    service.UsersService
	GroupService    service.GroupService
	MentorService   service.MentorService
	Config          config.Config
	UserRepository  repository.UsersRepository
}

func New(
	opts FxOpts,
) *Bot {
	return &Bot{
		DB:              opts.DB,
		Api:             opts.Api,
		MinioRepository: opts.MinioRepository,
		UsersService:    opts.UsersService,
		StudentService:  opts.StudentService,
		MentorService:   opts.MentorService,
		GroupService:    opts.GroupService,
		CacheRepository: opts.CacheRepository,
		CryptoKey:       []byte(opts.Config.CryptoKey),
		UserRepository:  opts.UserRepository,
	}
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	updates := b.Api.GetUpdatesChan(u)
	for update := range updates {
		go func(bot *Bot, update tgbotapi.Update) {
			ID := GetChatID(update)
			if ID != 0 {
				stack := userRuns[ID]
				if stack.Action != nil {
					stack.Update = &update
					userRuns[ID] = userRuns[ID].Action(stack)
				} else {
					if update.Message != nil {
						userDatas[ID] = &Data{LastMes: -1}
						userRuns[ID] = MainMenu(CallStack{
							ChatID:  ID,
							Bot:     bot,
							Update:  &update,
							IsPrint: true,
							Data:    update.Message.From.UserName,
						})
					}
				}
			}
		}(b, update)
	}
	return nil
}

func Chop(stack CallStack) CallStack {
	// Send "Work in progress"
	msg := tgbotapi.NewMessage(stack.ChatID, "Work in progress")
	_, err := stack.Bot.Api.Send(msg)
	if err != nil {
		log.Println(err)
	}
	// return on parent Run
	return ReturnOnParent(stack)
}

func GetChatID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}

	return -1
}

func ReturnOnParent(stack CallStack) CallStack {
	if stack.Parent != nil {
		if stack.Data == "Created3" {
			stack.Parent.Parent.Parent.IsPrint = true
			stack.Parent.Parent.Parent.Update = nil
			return stack.Parent.Parent.Parent.Action(*stack.Parent.Parent.Parent)
		}
		if stack.Data == "Created" {
			stack.Parent.Parent.IsPrint = true
			stack.Parent.Parent.Update = nil
			return stack.Parent.Parent.Action(*stack.Parent.Parent)
		}
		stack.Parent.IsPrint = true
		stack.Parent.Update = nil
		return stack.Parent.Action(*stack.Parent)
	}
	return RunTemplate(CallStack{
		IsPrint: true,
		ChatID:  stack.ChatID,
		Bot:     stack.Bot,
	})
}
