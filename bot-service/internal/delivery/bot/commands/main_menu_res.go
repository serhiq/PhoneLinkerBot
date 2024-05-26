package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const ADD_ACCOUNT_MESSAGE_ERROR = "Возникла ошибка при сохранении номера телефона"

func KeyboardMain() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(HELP_BUTTON),
		),
	)
}
