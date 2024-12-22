package handlers

import (
	"encoding/json"
	"go_final/app/internal/db"
	"go_final/app/internal/models"
	"go_final/app/internal/repeatTask"
	"net/http"
	"time"
)

func PostTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)

	var task = models.Remind{}
	// декодируем JSON
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// начинаем проверку даты
	now := time.Now()
	today := now.Format("20060102")
	// если поле date не указано - берется сегодняшнее число
	if task.Date == "" || task.Date == "today" {
		task.Date = today
	} else {
		// парсим дату задачи
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}
		// если дата меньше сегодняшнего числа, то
		if parsedDate.Before(time.Now()) {
			// правило не указано - дата сегодняшнее число
			if task.Repeat == "" {
				task.Date = today
				// правило указано - дата - та, что в правиле
			} else {
				nextDate, err := repeatTask.NextDate(time.Now(), task.Date, task.Repeat)
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

	// проверяем правило повторения в любом случае
	if task.Repeat != "" {
		_, err := repeatTask.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}
	id, err := db.InsertTask(task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка записи в БД"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}
