package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AsaHero/dastyor-bot/delivery/api"
	"github.com/AsaHero/dastyor-bot/delivery/telegram"
	"github.com/AsaHero/dastyor-bot/delivery/telegram/handler"
	"github.com/AsaHero/dastyor-bot/internal/inerr"
	llm_api "github.com/AsaHero/dastyor-bot/internal/infrastructure/llm"
	users_repo "github.com/AsaHero/dastyor-bot/internal/infrastructure/repository/users"
	"github.com/AsaHero/dastyor-bot/internal/usecase/llm"
	"github.com/AsaHero/dastyor-bot/internal/usecase/sessions"
	"github.com/AsaHero/dastyor-bot/internal/usecase/users"
	"github.com/AsaHero/dastyor-bot/pkg/config"
	"github.com/AsaHero/dastyor-bot/pkg/database/postgres"
	"github.com/AsaHero/dastyor-bot/pkg/logger"
	"github.com/AsaHero/dastyor-bot/pkg/redis"
	telegram_bot "github.com/AsaHero/dastyor-bot/pkg/telegram-bot"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	cfg         *config.Config
	logger      *logrus.Logger
	server      *http.Server
	telegramBot *telegram_bot.TelegramBot
	db          *gorm.DB
	redis       *redis.RedisStorage
}

func New(cfg *config.Config) *App {
	logger := logger.Init(cfg, cfg.APP+".log")

	db, err := postgres.NewGORMPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	redis, err := redis.NewRedisStorage(cfg)
	if err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}

	telegramBot := telegram_bot.NewTelegramBot(cfg.Bot.Token, cfg.Bot.WebhookURL)

	return &App{
		cfg:         cfg,
		logger:      logger,
		db:          db,
		telegramBot: telegramBot,
		redis:       redis,
	}
}

func (a *App) Run() error {
	// parse context timeout
	contextDeadline, err := time.ParseDuration(a.cfg.Context.Timeout)
	if err != nil {
		return inerr.WithMessage(err, "failed to parse context timeout")
	}

	// init llm api
	llmAPI, err := llm_api.New(a.cfg)
	if err != nil {
		return inerr.WithMessage(err, "failed to init llmAPI")
	}

	// init repo
	usersRepo := users_repo.New(a.db)

	// init services
	usersService := users.New(contextDeadline, usersRepo)
	sessionsService := sessions.New(contextDeadline, a.redis)
	llmService := llm.New(contextDeadline, llmAPI)

	// init bot router
	_ = telegram.NewRouter(a.cfg, a.telegramBot, &handler.Options{
		UserService:     usersService,
		SessionsService: sessionsService,
		LlmService:      llmService,
	})

	// init api router
	apiRouter := api.NewRouter(api.HandlerOptions{
		Config:      a.cfg,
		TelegramBot: a.telegramBot,
	})

	a.server, err = api.NewServer(a.cfg, apiRouter)
	if err != nil {
		return inerr.WithMessage(err, "failed to init server")
	}

	go a.telegramBot.Start()

	fmt.Println("Listen: ", "address", a.cfg.Server.Host+a.cfg.Server.Port)
	return a.server.ListenAndServe()
}

func (a *App) Stop() {
	a.server.Close()

	a.telegramBot.Stop()

	sqlDB, _ := a.db.DB()

	sqlDB.Close()

	a.logger.Writer().Close()
}
