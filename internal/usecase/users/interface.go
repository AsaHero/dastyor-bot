package users

import (
	"context"

	"github.com/AsaHero/dastyor-bot/internal/entity"
)

type Users interface {
	GetByExternalID(ctx context.Context, externalID int64) (*entity.Users, error)
}
