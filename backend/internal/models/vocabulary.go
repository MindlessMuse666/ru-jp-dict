package models

import "time"

/* Модель слова из словаря */
type Vocabulary struct {
	ID        int       `json:"id" db:"id"`
	Russian   string    `json:"russian" db:"russian"`
	Japanese  string    `json:"japanese" db:"japanese"`
	Onyomi    string    `json:"onyomi" db:"onyomi"`
	Kunyomi   string    `json:"kunyomi" db:"kunyomi"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
