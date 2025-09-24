package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Storage struct {
	DB *sql.DB
}

func NewDB(connStr string) (*Storage, error) {
	var db *sql.DB
	var err error

	const attempts = 20
	const wait = 2 * time.Second

	for i := 0; i < attempts; i++ {
		db, err = sql.Open("pgx", connStr)
		if err != nil {
			fmt.Printf("БД open failed (попытка %d/%d): %v\n", i+1, attempts, err)
			time.Sleep(wait)
			continue
		}
		if err = db.Ping(); err != nil {
			fmt.Printf("БД пинг неудачен (попытка %d/%d): %v\n", i+1, attempts, err)
			_ = db.Close()
			time.Sleep(wait)
			continue
		}
		fmt.Println("подключились к БД")
		return &Storage{DB: db}, nil
	}

	if err == nil {
		err = fmt.Errorf("неизвестная ошибка подключения к БД")
	}
	return nil, fmt.Errorf("Ошибка подключения к БД: %v", err)
}
