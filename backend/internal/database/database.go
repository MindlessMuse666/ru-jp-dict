package database

import (
	"database/sql"
	"log"

	"github.com/glebarez/sqlite"
)

/* Подключение к SQLite и создание таблиц, если они не существуют */
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(sqlite.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS vocabulary (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		russian TEXT NOT NULL,
		japanese TEXT NOT NULL,
		onyomi TEXT DEFAULT '',
		kunyomi TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("база данных инициализирована успешно")
	return db, nil
}
