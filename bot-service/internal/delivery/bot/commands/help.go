package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/PhoneLinkerBot/internal/app"
)

func AboutCommand(app *app.App, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, HELP_MESSAGE)
	return app.Bot.Reply(msg)
}
