package service

import (
	"context"

	"github.com/didsqq/user_api/internal/domain"
	"github.com/didsqq/user_api/internal/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, user domain.User) error {
	err := s.repo.UserRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
func (s *UserService) GetByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := s.repo.UserRepository.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
func (s *UserService) Update(ctx context.Context, user domain.UpdateUserInput) error {
	err := s.repo.UserRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
func (s *UserService) Delete(ctx context.Context, id int64) error {
	err := s.repo.UserRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
func (s *UserService) List(ctx context.Context, c domain.Conditions) ([]domain.User, error) {
	users, err := s.repo.UserRepository.List(ctx, c)
	if err != nil {
		return nil, err
	}

	return users, nil
}
