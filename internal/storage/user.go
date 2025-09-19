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
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1`
	u := User{}
	err := s.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Пользователь не найден: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("Ошибка поиска пользователя: %v", err)
	}
	return &u, nil
}
