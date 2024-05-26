package server

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	config "github.com/serhiq/PhoneLinkerBot/internal/configs"
	"github.com/serhiq/PhoneLinkerBot/internal/delivery/bot"
	repo "github.com/serhiq/PhoneLinkerBot/internal/repository/userPhoneMapping"
	"github.com/serhiq/PhoneLinkerBot/pkg/store/mysql"
	"log"
	"time"
)

const (
	longPollTimeout = 15
)

type Server struct {
	cfg               config.Config
	store             *mysql.Store
	delivery          *bot.Performer
	startFunc         []func()
	stopFunc          []func()
	sessionRepository *repo.Repository
}

func Serve(cfg config.Config) (*Server, error) {
	fmt.Printf("\n%#v\n\n", cfg)

	var s = &Server{
		cfg:       cfg,
		store:     nil,
		delivery:  nil,
		startFunc: nil,
		stopFunc:  nil,
	}

	for _, init := range []func() error{
		s.initDb,
		s.initSessionRepository,
		s.initBot,
	} {
		if err := init(); err != nil {
			return nil, errors.Wrap(err, "serve failed")
		}
	}
	return s, nil
}

func (s *Server) Start() error {
	fmt.Println("Server is starting...")

	for _, start := range s.startFunc {
		start()
	}

	return nil
}

func (s *Server) Stop() {
	for _, stop := range s.stopFunc {
		stop()
	}
}

func (s *Server) initSessionRepository() error {

	s.sessionRepository = repo.New(s.store.Db)
	return nil
}

func (s *Server) initDb() error {
	store, err := mysql.New(s.cfg.DBConfig)

	if err != nil {
		return err
	}

	s.store = store

	s.addStopDelegate(func() {
		db, err := s.store.Db.DB()
		if err != nil {
			log.Printf("database: error close database, %s", err)
			return
		}
		err = db.Close()
		if err != nil {
			log.Printf("database: error close database, %s", err)
			return
		}
		log.Print("database: close")
	})
	return err
}

func (s *Server) addStartDelegate(delegate func()) {
	s.startFunc = append(s.startFunc, delegate)
}

func (s *Server) addStopDelegate(delegate func()) {
	s.stopFunc = append(s.stopFunc, delegate)
}

func (s *Server) initBot() error {
	sBot, err := bot.New(bot.Options{
		Token: s.cfg.Telegram.Token,
	}, s.sessionRepository)

	if err != nil {
		return errors.Wrap(err, "cannot initialize Bot")
	}

	s.delivery = sBot

	u := tgbotapi.NewUpdate(0)
	u.Timeout = longPollTimeout

	updates := s.delivery.App.Bot.Api.GetUpdatesChan(u)
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	s.addStartDelegate(func() {
		//logger.SugaredLogger.Infof("Bot online %s", s.delivery.App.Bot.Api.Self.UserName)
		for update := range updates {
			go s.delivery.Dispatch(&update)
		}
	})

	s.addStopDelegate(func() {
		//logger.SugaredLogger.Info("Bot is stopping...")
		s.delivery.App.Bot.Api.StopReceivingUpdates()
	})

	return nil
}
