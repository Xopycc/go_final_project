package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"go-final-project/internal/api"
	"go-final-project/internal/config"
	"go-final-project/internal/repository"
	"go-final-project/internal/sqlidb"
	"go-final-project/internal/task"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Используем значение из конфига для подключения к базе данных
	db, err := sqlidb.Open("sqlite", cfg.DB)
	if err != nil {
		log.Println(err)
		return
	}

	repo := repository.New(db)
	srv := task.NewService(repo)
	api := api.New(srv)
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	webDir := filepath.Join(curDir, "web/")
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	// Исправление маршрутов без HTTP метода (GET, POST, и т.д.)
	mux.HandleFunc("/api/nextdate", api.NextDate)
	mux.HandleFunc("/api/task", api.TaskCreate)
	mux.HandleFunc("/api/tasks", api.GetTasks)
	mux.HandleFunc("/api/task", api.GetTask)
	mux.HandleFunc("/api/task", api.UpdateTask)
	mux.HandleFunc("/api/task/done", api.TaskDone)
	mux.HandleFunc("/api/task", api.TaskDelete)

	log.Printf("Сервер запущен на порту %s\n", cfg.Port)

	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
