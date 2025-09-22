package service

import (
	"context"
	"errors"

	"github.com/Royal17x/go-url-shortener/internal/storage"
	"github.com/Royal17x/go-url-shortener/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Store *storage.Storage
}

func NewAuthService(store *storage.Storage) *AuthService {
	return &AuthService{Store: store}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (int, error) {
	if username == "" || email == "" || password == "" {
		return 0, errors.New("username, email, password обязательны")
	}

	if !validation.ValidateEmail(email) {
		return 0, errors.New("некорректный email")
	}

	if len(password) < 6 {
		return 0, errors.New("пароль должен быть от 6 символов и больше")
	}

	existing, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, err
	}

	if existing != nil {
		return 0, errors.New("пользователь уже существует")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("ошибка при хэшировании пароля")
	}

	uid, err := s.Store.CreateUser(ctx, storage.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	})

	if err != nil {
		return 0, err
	}

	return uid, nil
}
