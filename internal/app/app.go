package app

import (
	"log"
	"time"

	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"github.com/AsaHero/dastyor-bot/pkg/config"
	"github.com/AsaHero/dastyor-bot/pkg/database/postgres"
	"github.com/AsaHero/dastyor-bot/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	logger *logrus.Logger
	db     *gorm.DB
}

func New(cfg *config.Config) *App {
	logger := logger.Init(cfg, cfg.APP+".log")

	db, err := postgres.NewGORMPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	return &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

func (a *App) Run() error {
	// parse context timeout
	_, err := time.ParseDuration(a.cfg.Context.Timeout)
	if err != nil {
		return inerr.Err(err)
	}

	return nil
}

func (a *App) Stop() {
}
