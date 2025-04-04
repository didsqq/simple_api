package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/didsqq/user_api/internal/domain"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		db: db,
	}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user domain.User) error {
	const op = "UserRepository.Create"

	createQuery := fmt.Sprintf("INSERT INTO %s (name, email) VALUES ($1, $2)", usersTable)

	stmt, err := r.db.PrepareContext(ctx, createQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepositoryPostgres) GetByID(ctx context.Context, id int64) (domain.User, error) {
	const op = "UserRepository.GetByID"

	user := domain.User{}

	getQuery := fmt.Sprintf("SELECT id, name, email FROM %s WHERE id=$1 AND deleted_at IS NULL", usersTable)

	stmt, err := r.db.PrepareContext(ctx, getQuery)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: user with ID %d not found", op, id)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepositoryPostgres) Update(ctx context.Context, user domain.UpdateUserInput) error {
	const op = "UserRepository.Update"

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *user.Name)
		argId++
	}

	if user.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *user.Email)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", usersTable, setQuery, argId)

	args = append(args, user.ID)

	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (r *UserRepositoryPostgres) Delete(ctx context.Context, id int64) error {
	const op = "UserRepository.Delete"

	deleteQuery := fmt.Sprintf("UPDATE %s SET deleted_at=$1 WHERE id=$2", usersTable)
	stmt, err := r.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, time.Now(), id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (r *UserRepositoryPostgres) List(ctx context.Context, c domain.Conditions) ([]domain.User, error) {
	const op = "UserRepository.List"

	selectQuery := fmt.Sprintf("SELECT id, name, email FROM %s WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", usersTable)
	stmt, err := r.db.PrepareContext(ctx, selectQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, c.Limit, c.Offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}
