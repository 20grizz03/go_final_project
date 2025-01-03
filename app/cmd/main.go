package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go_final/app/internal/handlers"
	"go_final/app/pkg/config"
	"log"
	"net/http"
	"os"
)

var (
	portLocal = "7540"
	webDir    = "./web"
)

func main() {
	config.LoadEnviroment()
	config.MakeDB()
	defer config.CloseDB()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = portLocal
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	// вычисление следующей даты
	r.Get("/api/nextdate", handlers.HandlerForNewDate)
	// добавление следующей задачи
	r.Post("/api/task", handlers.PostTask)
	// отображение задач
	r.Get("/api/tasks", handlers.GetTasks)
	// получение одной задачи
	r.Get("/api/task", handlers.GetTask)
	// изменение одной задачи
	r.Put("/api/task", handlers.PutTask)
	// выполнение задачи
	r.Post("/api/task/done", handlers.DoneTask)
	// удаление задачи
	r.Delete("/api/task", handlers.DeleteTask)

	// запуск сервера
	log.Printf("Сервер запущен на http://localhost:%s/", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
