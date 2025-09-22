package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

func (s *Storage) CreateUser(ctx context.Context, u User) (int, error) {
	query := `INSERT INTO users (username, email, password_hash)
			  VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := s.DB.QueryRowContext(ctx, query, u.Username, u.Email, u.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Ошибка создания пользователя: %v", err)
	}
	return id, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.DB.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash FROM users WHERE email = $1`, email)
	u := User{}
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	return &u, nil
}
