package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/models"
)

type VocabularyRepo struct {
	db *sql.DB
}

func NewVocabularyRepo(db *sql.DB) *VocabularyRepo {
	return &VocabularyRepo{db: db}
}

// GET Получает слово по ID
func (r *VocabularyRepo) GetByID(id int) (models.Vocabulary, error) {
	query := `
	SELECT
		id, russian, japanese, onyomi, kunyomi, created_at, updated_at
	FROM
        vocabulary
    WHERE
        id = ?
	`

	var word models.Vocabulary
	err := r.db.QueryRow(query, id).Scan(
		&word.ID,
		&word.Russian,
		&word.Japanese,
		&word.Onyomi,
		&word.Kunyomi,
		&word.CreatedAt,
		&word.UpdatedAt,
	)

	if err != nil {
		return models.Vocabulary{}, err
	}

	return word, nil
}

// GET Получает все слова
func (r *VocabularyRepo) GetAll() ([]models.Vocabulary, error) {
	query := `
	SELECT
		*
	FROM
		vocabulary
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.Vocabulary

	for rows.Next() {
		var w models.Vocabulary

		err := rows.Scan(
			&w.ID,
			&w.Russian,
			&w.Japanese,
			&w.Onyomi,
			&w.Kunyomi,
			&w.CreatedAt,
			&w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// POST Создает новое слово
func (r *VocabularyRepo) Create(v models.Vocabulary) (int, error) {
	query := `
	INSERT INTO
		vocabulary (russian, japanese, onyomi, kunyomi)
	VALUES
		(?, ?, ?, ?);
	`

	result, err := r.db.Exec(query, v.Russian, v.Japanese, v.Onyomi, v.Kunyomi)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// PATCH Частично обновляет слово по ID
func (r *VocabularyRepo) PartialUpdate(id int, updates map[string]any) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Начало построения SQL-запроса
	query := "UPDATE vocabulary SET "
	var setClauses []string
	var values []interface{}

	// Добавление полей для обновления
	if russian, ok := updates["russian"]; ok {
		setClauses = append(setClauses, "russian = ?")
		values = append(values, russian)
	}
	if japanese, ok := updates["japanese"]; ok {
		setClauses = append(setClauses, "japanese = ?")
		values = append(values, japanese)
	}
	if onyomi, ok := updates["onyomi"]; ok {
		setClauses = append(setClauses, "onyomi = ?")
		values = append(values, onyomi)
	}
	if kunyomi, ok := updates["kunyomi"]; ok {
		setClauses = append(setClauses, "kunyomi = ?")
		values = append(values, kunyomi)
	}

	// Обновление updated_at
	setClauses = append(setClauses, "updated_at = ?")
	values = append(values, time.Now())

	// Добавление условия WHERE
	query += strings.Join(setClauses, ", ") + " WHERE id = ?"
	values = append(values, id)

	// Выполнение SQL-запроса
	_, err := r.db.Exec(query, values...)
	return err
}

// PUT Полностью обновляет слово по ID
func (r *VocabularyRepo) Update(id int, v models.Vocabulary) error {
	query := `
	UPDATE vocabulary
	SET
		russian = ?,
		japanese = ?,
		onyomi = ?,
		kunyomi = ?,
		updated_at = ?
	WHERE
		id = ?
	`

	_, err := r.db.Exec(
		query,
		v.Russian,
		v.Japanese,
		v.Onyomi,
		v.Kunyomi,
		time.Now(),
		id,
	)
	return err
}

// DELETE Удаляет слово по ID
func (r *VocabularyRepo) Delete(id int) error {
	query := `
	DELETE FROM vocabulary
	WHERE
		id = ?
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	// Проверка, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Выбрасываем ошибку, если ни одна строка не была удалена
	if rowsAffected == 0 {
		return fmt.Errorf("слово с id=%d не найдено", id)
	}

	return nil
}
