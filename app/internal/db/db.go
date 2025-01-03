package db

import (
	"database/sql"
	"errors"
	"go_final/app/internal/models"
	"go_final/app/pkg/config"
	"time"
)

func InsertInDB(task models.Remind) (uint64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := config.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	// получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// реализуем поиск в БД
func FindInDb(search string, limit int) ([]models.Remind, error) {
	var query string
	var args []interface{}

	// Если параметр search пустой, возвращаем все задачи
	if search == "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
		args = append(args, limit)
	} else {
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			// если дата, используем сравнение по дате
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?"
			args = append(args, parsedDate.Format("20060102"), limit)
		} else {
			// если подстрока - в title и comment
			likePattern := "%" + search + "%"
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			args = append(args, likePattern, likePattern, limit)
		}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, errors.New("ошибка выполнения запроса к базе данных")
	}
	defer rows.Close()

	var tasks []models.Remind
	for rows.Next() {
		var task models.Remind
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, errors.New("ошибка сканирования данных из базы")
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("ошибка постобработки данных из базы")
	}
	return tasks, nil
}

// получение задачи по ID
func GetTaskByID(id int) (models.Remind, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := config.DB.QueryRow(query, id)
	var task models.Remind
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return models.Remind{}, errors.New("Задача не найдена")
	}
	return task, nil
}

// обновление функции
func UpdateTask(task models.Remind) error {
	query := `
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?`
	res, err := config.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// удаление из БД
func DeleteTaskByID(id int) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := config.DB.Exec(query, id)
	return err
}

// обновление задачи только при наличии remind по дате
func UpdateTaskDate(id uint64, newDate string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	_, err := config.DB.Exec(query, newDate, id)
	return err
}
