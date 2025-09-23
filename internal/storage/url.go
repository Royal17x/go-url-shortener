package storage

import (
	"context"
	"fmt"
)

type URL struct {
	ID          int
	UserID      int
	ShortCode   string
	OriginalURL string
}

func (s *Storage) CreateURL(ctx context.Context, u URL) (int, error) {
	query := `INSERT INTO urls (user_id, short_code, original_url)
			  VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.DB.QueryRowContext(ctx, query, u.UserID, u.ShortCode, u.OriginalURL)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания ссылки:%v", err)
	}
	return id, nil
}
