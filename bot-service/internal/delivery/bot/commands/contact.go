package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/serhiq/PhoneLinkerBot/internal/app"
	"github.com/serhiq/PhoneLinkerBot/internal/repository/models"
	"gorm.io/gorm"
)

func ContactHandler(app *app.App, phone string, message *tgbotapi.Message) error {
	_, err := app.Repository.GetByChatID(message.Chat.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return handleNewMapping(app, message.Chat.ID, phone)
		}
		return err
	}

	// Если номер телефона уже существует, то смыла обновлять его смысла нет
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ваш номер телефона уже сохранён.")
	return app.Bot.Reply(msg)
}

func handleNewMapping(app *app.App, chatID int64, phone string) error {
	mapping := models.NewUserPhoneMapping(chatID, phone)
	err := app.Repository.InsertOrUpdate(mapping)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, ADD_ACCOUNT_MESSAGE_ERROR)
		return app.Bot.Reply(errorMsg)
	}

	msg := tgbotapi.NewMessage(chatID, "Ваш номер телефона успешно сохранён.")
	return app.Bot.Reply(msg)
}
