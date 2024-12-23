package handlers

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"go_final/app/internal/db"
	"go_final/app/internal/repeatTask"
	"net/http"
	"strconv"
	"time"
)

func DoneTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id") // Извлечение параметра из пути
	if idStr == "" {               // Если параметр пустой, пробуем из строки запроса
		idStr = r.URL.Query().Get("id")
	}

	if idStr == "" { // Если параметр всё ещё пустой, возвращаем ошибку
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"Некорректный идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// получаем задачу из базы данных
	task, err := db.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Ошибка получения задачи из БД"}`, http.StatusInternalServerError)
		}
		return
	}

	// проверка на пустоту продолжения задачи
	if task.Repeat == "" {
		err := db.DeleteTaskByID(id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := repeatTask.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка расчёта следующей даты повторения"}`, http.StatusInternalServerError)
			return
		}

		// Обновляем дату задачи в базе данных
		err = db.UpdateTaskDate(uint64(id), nextDate)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}
