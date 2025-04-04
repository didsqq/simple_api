package validate

import (
	"errors"

	"github.com/didsqq/user_api/internal/domain"
)

func ValidateUser(user domain.User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func ValidateUpdateUser(user domain.UpdateUserInput) error {
	if user.Email == nil && user.Name == nil {
		return errors.New("no updates")
	}

	return nil
}
