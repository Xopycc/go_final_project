package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	mux.HandleFunc("GET /api/nextdate", api.NextDate)
	mux.HandleFunc("POST /api/task", api.TaskCreate)
	mux.HandleFunc("GET /api/tasks", api.GetTasks)
	mux.HandleFunc("GET /api/task", api.GetTask)
	mux.HandleFunc("PUT /api/task", api.UpdateTask)
	mux.HandleFunc("POST /api/task/done", api.TaskDone)
	mux.HandleFunc("DELETE /api/task", api.TaskDelete)

	// Добавляем двоеточие перед портом, если оно отсутствует
	port := cfg.Port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("Сервер запущен на порту %s\n", port)

	log.Fatal(http.ListenAndServe(port, mux))
}
