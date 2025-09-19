package storage

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	DB *sql.DB
}

func NewDB(connStr string) (*Storage, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("БД не отвечает: %v", err)
	}
	return &Storage{DB: db}, nil
}
