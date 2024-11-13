package users

import (
	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/infrastructure/repository"
)

type Repository interface {
	repository.BaseRepository[entity.Users]
}
