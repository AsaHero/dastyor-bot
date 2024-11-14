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

func New(contextTimeout time.Duration, repo users.Repository) Users {
	return &usersService{
		contextTimeout: contextTimeout,
		repo:           repo,
	}
}

func (s *usersService) GetByExternalID(ctx context.Context, userID int64) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repo.FindOne(ctx, map[string]any{"id": userID})
	if err != nil {
		return nil, inerr.Err(err)
	}

	return user, nil
}

func (s *usersService) Upsert(ctx context.Context, user *entity.Users) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err := s.repo.Upsert(ctx, []string{"username", "first_name", "last_name"}, user)
	if err != nil {
		return inerr.Err(err)
	}

	return nil
}
