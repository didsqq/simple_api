package repository

import (
	"context"
	"database/sql"

	"github.com/didsqq/user_api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, user domain.UpdateUserInput) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, c domain.Conditions) ([]domain.User, error)
}

type Repository struct {
	UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepositoryPostgres(db),
	}
}
