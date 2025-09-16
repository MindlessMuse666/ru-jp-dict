package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/models"
)

type VocabularyRepo struct {
	db *sql.DB
}

func NewVocabularyRepo(db *sql.DB) *VocabularyRepo {
	return &VocabularyRepo{db: db}
}

/* Получить все слова из БД */
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

/* Добавить новое слово в БД */
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

/* Обновляет существующее слово по ID */
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

/* Удаляет слово по ID */
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
