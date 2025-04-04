package repository

import (
	"database/sql"
	"fmt"
)

const (
	usersTable = "users"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	const op = "repository.NewPostgresDB"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: ошибка при подключении к БД: %w", op, err)
	}

	return db, nil
}
