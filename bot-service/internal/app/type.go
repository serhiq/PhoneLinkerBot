package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/PhoneLinkerBot/internal/repository"
)

type App struct {
	Repository repository.UserPhoneMappingRepository
	Cfg        *AppConfig

	Bot *TelegramBot
}

type AppConfig struct {
}

type TelegramBot struct {
	Api *tgbotapi.BotAPI
}

func NewTelegramBot(b *tgbotapi.BotAPI) *TelegramBot {
	return &TelegramBot{Api: b}
}

func (t TelegramBot) Reply(msg tgbotapi.Chattable) error {
	_, err := t.Api.Send(msg)

	if err != nil {
		return NewErrorRespond(err)
	}

	return nil
}

/*
send an empty callback response for prevent the "waiting" icon from appearing on an inline keyboard
*/
func (t TelegramBot) AnswerEmptyOnCallback(callbackID string) error {
	answer := tgbotapi.CallbackConfig{
		CallbackQueryID: callbackID,
	}
	return t.Reply(answer)

}

func (t TelegramBot) ReplyWithId(msg tgbotapi.Chattable) (*tgbotapi.Message, error) {
	resultMsg, err := t.Api.Send(msg)

	if err != nil {
		return nil, NewErrorRespond(err)
	}

	return &resultMsg, nil
}

func NewErrorRespond(err error) *ErrRespond {
	return &ErrRespond{
		err: err.Error(),
	}

}

type ErrRespond struct {
	err string
}

func (e ErrRespond) Error() string {
	return fmt.Sprintf("Failed to respond  %s", e.err)
}

func IsRespondError(err error) bool {
	_, ok := err.(ErrRespond)
	return ok
}
