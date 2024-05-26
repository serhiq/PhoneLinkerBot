package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/serhiq/PhoneLinkerBot/internal/app"
	"gorm.io/gorm"
	"log"
)

/*
   /start
*/

func StartCommand(app *app.App, message *tgbotapi.Message) error {
	log.Printf("Executing StartCommand for chat ID: %d", message.Chat.ID)

	msg := tgbotapi.NewMessage(message.Chat.ID, START_MESSAGE)
	_, err := app.Repository.GetByChatID(message.Chat.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Chat ID: %d not found, requesting contact", message.Chat.ID)
			requestContact := tgbotapi.NewMessage(message.Chat.ID, REQUEST_CONTACT_PHONE_MESSAGE)
			requestContact.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonContact(SEND_PHONE_BUTTON),
			))

			return app.Bot.Reply(requestContact)
		}
		return err
	}

	msg.ReplyMarkup = KeyboardMain()
	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}
	return nil

}
