package service

import (
	"context"

	"github.com/didsqq/user_api/internal/domain"
	"github.com/didsqq/user_api/internal/repository"
)

type User interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, user domain.UpdateUserInput) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, c domain.Conditions) ([]domain.User, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo),
	}
}
