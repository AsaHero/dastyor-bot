package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/inerr"
	"github.com/AsaHero/dastyor-bot/pkg/redis"
)

type sessionsService struct {
	contextDeadline time.Duration
	cache           *redis.RedisStorage
}

func New(contextDeadline time.Duration, cache *redis.RedisStorage) Sessions {
	return &sessionsService{
		contextDeadline: contextDeadline,
		cache:           cache,
	}
}

func (sessionsService) prefixKey(userId int64) string {
	return fmt.Sprintf("session:%d", userId)
}

func (s *sessionsService) Set(ctx context.Context, session *entity.Sessions) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextDeadline)
	defer cancel()

	err := s.cache.Save(ctx, s.prefixKey(session.UserID), session)
	if err != nil {
		return inerr.Err(err)
	}

	return nil
}

func (s *sessionsService) Get(ctx context.Context, userID int64) (*entity.Sessions, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextDeadline)
	defer cancel()

	session := &entity.Sessions{}

	err := s.cache.Get(ctx, s.prefixKey(userID), session)
	if err != nil {
		return nil, inerr.Err(err)
	}

	return session, nil
}
