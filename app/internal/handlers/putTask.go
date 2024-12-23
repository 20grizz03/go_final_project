package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final/app/internal/db"
	"go_final/app/internal/models"
	"go_final/app/internal/repeatTask"
	"net/http"
	"time"
)

func PutTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var task models.Remind
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	// проверяем наличие поля id
	if task.ID == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// проверяем также поле title
	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// проверяем дату
	now := time.Now()
	today := now.Format("20060102")
	if task.Date == "" || task.Date == "today" {
		task.Date = today
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(time.Now()) {
			if task.Repeat == "" {
				task.Repeat = today
			} else {
				nextDate, err := repeatTask.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
					return
				}
				if task.Date != today {
					task.Date = nextDate
				}
			}
		}
	}
	// проверяем правило повторения
	if task.Repeat != "" {
		_, err := repeatTask.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}

	// Обновляем задачу в базе данных
	err := db.UpdateTask(task)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
		}
		return
	}
	// Возвращаем пустой JSON-объект
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))
}
