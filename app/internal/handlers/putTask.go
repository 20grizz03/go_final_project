package handlers

import (
	"encoding/json"
	"go_final/app/internal/models"
	"net/http"
)

func PutTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)

	var task = models.Remind{}

	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}
}
