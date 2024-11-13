package users

import (
	"github.com/AsaHero/dastyor-bot/internal/entity"
	"github.com/AsaHero/dastyor-bot/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type usersRepository struct {
	repository.BaseRepository[entity.Users]
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) Repository {
	return &usersRepository{
		BaseRepository: repository.NewBaseRepository[entity.Users](db),
		db:             db,
	}
}
