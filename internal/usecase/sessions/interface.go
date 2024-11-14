package sessions

import (
	"context"

	"github.com/AsaHero/dastyor-bot/internal/entity"
)

type Sessions interface {
	Get(ctx context.Context, userID int64) (*entity.Sessions, error)
	Set(ctx context.Context, session *entity.Sessions) error
}
