package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/PhoneLinkerBot/internal/app"
	"github.com/serhiq/PhoneLinkerBot/internal/delivery/bot/commands"
	"github.com/serhiq/PhoneLinkerBot/internal/repository"
)

type Options struct {
	Token string
}

type Performer struct {
	App            *app.App
	commandHandler []commandHandler
	contactHandler func(app *app.App, phone string, message *tgbotapi.Message) error
}

type commandHandler struct {
	name    string
	handler func(app *app.App, message *tgbotapi.Message) error
}

func (p *Performer) Dispatch(update *tgbotapi.Update) {
	if err := p.process(update); err != nil {
		p.processError(err, *update)
	}
}

func (p *Performer) isCommandButton(text string) bool {
	for _, hanler := range p.commandHandler {
		if hanler.name == text {
			return true
		}
	}
	return false
}

func (p *Performer) AddCommandHandler(name string, handl func(app *app.App, message *tgbotapi.Message) error) {
	ch := commandHandler{
		name:    name,
		handler: handl,
	}
	commandHanlers := append(p.commandHandler, ch)
	p.commandHandler = commandHanlers
}

func (p *Performer) RegisterContactHandler(handl func(app *app.App, phone string, message *tgbotapi.Message) error) {
	p.contactHandler = handl
}

func (p *Performer) processCommand(message *tgbotapi.Message) error {

	for _, command := range p.commandHandler {
		if command.name == message.Text {
			return command.handler(p.App, message)
		}
	}
	return commands.NewCommandNotFound(message.Text)
}

func (p *Performer) process(update *tgbotapi.Update) error {

	if update.Message != nil && p.isCommandButton(update.Message.Text) {
		return p.processCommand(update.Message)
	}

	if update.Message != nil && update.Message.Contact != nil && update.Message.Contact.PhoneNumber != "" {

		return p.contactHandler(p.App, update.Message.Contact.PhoneNumber, update.Message)
	}

	if update.Message.Text != "" {
		return commands.NewCommandNotFound(update.Message.Text)
	}

	return nil
}

func (p *Performer) processError(err error, update tgbotapi.Update) {

	if err != nil {
		if commands.IsCommandNotFoundError(err) {
			// todo если запрос неизвесен показать о приложении
			// если записан номер телефона  Показывать Наш сервис, предоставьте номер телефона

			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Простите, я не понял Ваш запрос.")
			message.ReplyMarkup = commands.KeyboardMain()

			err := p.App.Bot.Reply(message)
			if err != nil {
				return
			}

			//logger.SugaredLogger.Errorw("bot_command_not_found", "update", update,
			//	"chatId", update.FromChat(), "err", err)
			return
		}

		//logger.SugaredLogger.Errorw("bot_update",
		//	"chatId", update.FromChat(), "err", err)
		return
	}
}

func New(options Options, repo repository.UserPhoneMappingRepository) (*Performer, error) {
	bot, err := tgbotapi.NewBotAPI(options.Token)
	if err != nil {
		return nil, err
	}

	var p = Performer{
		App: &app.App{
			Repository: repo,
			Bot:        app.NewTelegramBot(bot),
			Cfg:        &app.AppConfig{},
		}}

	p.AddCommandHandler("/start", commands.StartCommand)
	p.AddCommandHandler("/help", commands.AboutCommand)
	p.AddCommandHandler(commands.HELP_BUTTON, commands.AboutCommand)
	p.RegisterContactHandler(commands.ContactHandler)

	bot.Debug = false

	return &p, nil
}
