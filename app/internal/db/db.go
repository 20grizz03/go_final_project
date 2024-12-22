package db

import (
	"go_final/app/internal/models"
	"go_final/app/pkg/config"
)

func InsertTask(task models.Remind) (uint64, error) {
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
