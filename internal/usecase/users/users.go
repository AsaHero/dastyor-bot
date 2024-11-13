package users

import (
	"context"
	"time"

	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"github.com/AsaHero/dastyor-bot/internal/infrastructure/repository/users"
)

type usersService struct {
	contextTimeout time.Duration
	repo           users.Repository
}

func NewUsersService(contextTimeout time.Duration, repo users.Repository) Users {
	return &usersService{
		contextTimeout: contextTimeout,
		repo:           repo,
	}
}

func (s *usersService) GetByExternalID(ctx context.Context, externalID int64) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repo.FindOne(ctx, map[string]any{"external_id": externalID})
	if err != nil {
		return nil, inerr.Err(err)
	}

	return user, nil
}
