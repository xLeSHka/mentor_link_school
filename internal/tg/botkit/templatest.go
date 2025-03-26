package botkit

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"strings"
)

var (
	MainMenuTemplate         = "[Начальное меню]"
	MainMenuTextTemplate     = "Пожалуйста выберите действие 🙏"
	AuthedMenuTemplate       = "[Главное меню]"
	RegisterMenuTemplate     = "[Регистрационное меню]"
	RegisterMenuTextTemplate = func(id uuid.UUID, name, telegram, password string) string {
		return fmt.Sprintf("Вы 🫵\nID: %s 🆔\nИмя: %s 🪪\nТелеграм: %s \nПароль: %s 🔑", id, name, telegram, password)
	}
	ValidatePasswordTemplate = func(minLength, maxLength, hasLower, hasUpper, hasDigit, hasSymbol bool) string {
		build := strings.Builder{}
		if minLength {
			build.Write([]byte(fmt.Sprintf("Минимальная длина: ✔️\n")))
		} else {
			build.Write([]byte(fmt.Sprintf("Минимальная длина: ❌\n")))
		}
		if maxLength {
			build.Write([]byte(fmt.Sprintf("Максимальная длина: ✔️\n")))
		} else {
			build.Write([]byte(fmt.Sprintf("Максимальная длина: ❌\n")))
		}
		if hasLower {
			build.Write([]byte(fmt.Sprintf("Строчная буква: ✔️\n")))
		} else {
			build.Write([]byte(fmt.Sprintf("Строчная буква: ❌\n")))
		}
		if hasUpper {
			build.Write([]byte(fmt.Sprintf("Прописная буква: ✔️\n")))
		} else {
			build.Write([]byte(fmt.Sprintf("Прописная буква: ❌\n")))
		}
		if hasDigit {
			build.Write([]byte(fmt.Sprintf("Цифра: ✔️\n")))
		} else {
			build.Write([]byte(fmt.Sprintf("Цифра: ❌\n")))
		}
		if hasSymbol {
			build.Write([]byte(fmt.Sprintf("Символ: ✔️\n")))
		} else {
			build.Write([]byte(fmt.Sprintf("Символ: ❌\n")))
		}
		return build.String()
	}
	LoginMenuTemplate         = "[Меню авторизации]"
	ErrorMenuTemplate         = "[Ошибка]"
	InternalErrorTextTemplate = "Приносим свои извинения, произошла непредвиденная ошибка! 🥺🙏\nВведите /start чтобы начать с начала!"
	ProfileTextTemplate       = func(id uuid.UUID, name, bio string) string {
		return fmt.Sprintf("Пользователь\nID: %s\nИмя: %s\nБИО: %s", id, name, bio)
	}
	CreateGroupMenuTemplate = "[Меню создания организации]"
	CreateGroupTextTemplate = func(id uuid.UUID, name, inviteCode string) string {
		return fmt.Sprintf("Организация\nID: %s\nИмя: %s\nПригласительный код: %s", id, name, inviteCode)
	}
	GroupMenuTemplate = "[Меню организации]"
	GroupTextTemplate = func(id uuid.UUID, name string, inviteCode *string) string {
		builder := strings.Builder{}
		builder.Write([]byte(fmt.Sprintf("Организация\nID: %s\nИмя: %s", id, name)))
		if inviteCode != nil {
			builder.Write([]byte(fmt.Sprintf("\nПригласительный код: %s", *inviteCode)))
		}
		return builder.String()
	}
	MembersMenuTemplate     = "[Меню членов организации]"
	MembersMenuTextTemplate = func() string {
		return fmt.Sprintf("Выберите члена организации!")
	}
	MemberMenuTemplate     = "[Меню члена организации]"
	MemberMenuTextTemplate = func(id uuid.UUID, name, telegram, role string, isPair bool, bio *string, roles []*models.Role) string {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("Пользователь\nID: %s\nИмя: %s\n", id, name))
		if isPair {
			builder.WriteString(fmt.Sprintf("Телеграм: @%s", telegram))
		}
		if bio != nil {
			builder.WriteString(fmt.Sprintf("\nБИО: %s", *bio))
		}
		if role == "owner" {
			builder.WriteString("\nРоли: ")
			for _, role := range roles {
				switch role.Role {
				case "owner":
					builder.WriteString("🧑‍💼")
				case "mentor":
					builder.WriteString("🧑‍🏫")
				case "student":
					builder.WriteString("👨‍🎓")
				}
			}
		}
		return builder.String()
	}

	JoinMenuTemplate     = "[Меню входа в организацию]"
	GroupsMenuTemplate   = "[Меню выбора организации]"
	EditUserMenuTemplate = "[Меню редактирования пользователя]"
	StatMenuTemplate     = "[Меню статистики]"
	StatMenuTextTemplate = func(stat *models.GroupStat) string {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("Колличество студентов: %d👨‍🎓\nКолличество менторов: %d🧑‍🏫\nВсего запросов: %d\nПринятых запросов: %d✔️\nОтклонённых запросов: %d❌\nКонверсия запросов: %.2f", stat.StudentsCount, stat.MentorsCount, stat.HelpRequestCount, stat.AcceptedRequestCount, stat.RejectedRequestCount, stat.Conversion))
		return builder.String()
	}
	StudentsMenuTemplate     = "[Меню студентов]"
	StudentsMenuTextTemplate = func() string {
		return fmt.Sprintf("Выберите студента!")
	}
	MentorsMenuTemplate     = "[Меню менторов]"
	MentorsMenuTextTemplate = func() string {
		return fmt.Sprintf("Выберите ментора!")
	}
	AvailableMentorsMenuTemplate     = "[Меню доступных менторов]"
	AvailableMentorsMenuTextTemplate = func() string {
		return fmt.Sprintf("Выберите ментора!")
	}
	SendReqMenuTemplate            = "[Меню отправки запроса]"
	SendedRequestsMenuTemplate     = "[Меню исходящих запросов]"
	SendedRequestsMenuTextTemplate = func() string {
		return fmt.Sprintf("Выберите запрос!")
	}
	RecievedRequestsMenuTemplate     = "[Меню входящих запросов]"
	RecievedRequestsMenuTextTemplate = func() string {
		return fmt.Sprintf("Выберите запрос!")
	}
	RequestMenuTemplate     = "[Меню запроса]"
	RequestMenuTextTemplate = func(req *models.HelpRequest) string {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("ID: %s\nСтудент: %s\nМентор: %s\nID организации: %s\nЦель: %s", req.ID, req.Student.Name, req.Mentor.Name, req.GroupID, req.Goal))
		if req.BIO != nil {
			builder.WriteString(fmt.Sprintf("\nБИО студента: %s", *req.BIO))
		}
		builder.WriteString(fmt.Sprintf("\nСтатус: %s", Status(req.Status)))
		return builder.String()
	}
	Status = func(status string) string {
		switch status {
		case "accepted":
			return "✔️"
		case "rejected":
			return "❌"
		default:
			return "⌛"
		}
	}
)
