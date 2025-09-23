package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type URL struct {
	ID          int
	UserID      int
	ShortCode   string
	OriginalURL string
}

func (s *Storage) SaveURL(ctx context.Context, u URL) (int, error) {
	query := `INSERT INTO urls (user_id, short_code, original_url)
			  VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.DB.QueryRowContext(ctx, query, u.UserID, u.ShortCode, u.OriginalURL)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания ссылки:%v", err)
	}
	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, shortCode string) (*URL, error) {
	query := `SELECT id, user_id, original_url, short_code FROM urls WHERE short_code = $1`
	u := URL{}
	err := s.DB.QueryRowContext(ctx, query, shortCode).Scan(&u.ID, &u.UserID, &u.OriginalURL, &u.ShortCode)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
