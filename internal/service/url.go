package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/Royal17x/go-url-shortener/internal/storage"
)

type URLService struct {
	Store *storage.Storage
}

func NewURLService(store *storage.Storage) *URLService {
	return &URLService{Store: store}
}

func generateShortCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

func (s *URLService) ShortenURL(ctx context.Context, userID int, originalURL string) (string, error) {
	code := generateShortCode()
	_, err := s.Store.SaveURL(ctx, storage.URL{
		UserID:      userID,
		ShortCode:   code,
		OriginalURL: originalURL,
	})
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *URLService) ResolveURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.Store.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}
