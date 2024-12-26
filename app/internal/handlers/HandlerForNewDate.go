package handlers

import (
	"go_final/app/internal/repeatTask"
	"net/http"
	"time"
)

func HandlerForNewDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, `{"error": "некорректная дата now"}`, http.StatusBadRequest)
		return
	}
	// Вызываем функцию NextDate
	nextDate, err := repeatTask.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(nextDate))
}
