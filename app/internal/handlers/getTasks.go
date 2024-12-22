package handlers

import (
	"encoding/json"
	"fmt"
	"go_final/app/internal/db"
	"net/http"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	// ограничение на кол-во возвращаемых задач
	const limit = 50

	tasks, err := db.FindInDB(search, limit)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// преобразуем задачи в формат, ожидаемый тестами
	var tasksForResponse []map[string]string
	for _, task := range tasks {
		tasksForResponse = append(tasksForResponse, map[string]string{
			"id":      fmt.Sprintf("%d", task.ID),
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		})
	}

	// если пустой, то возвращаем слайс
	if tasksForResponse == nil {
		tasksForResponse = []map[string]string{}
	}

	// формируем JSON-ответ
	response := map[string]interface{}{
		"tasks": tasksForResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Ошибка формирования JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
}
