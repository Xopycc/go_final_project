package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go-final-project/internal/task"
)

type errResponse struct {
	Error string `json:"error"`
}

func sendErr(w http.ResponseWriter, errText string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	rsp := errResponse{Error: errText}
	err := json.NewEncoder(w).Encode(rsp)
	if err != nil {
		log.Println(err)
	}
}

type idResponse struct {
	ID int `json:"id"`
}

func sendID(w http.ResponseWriter, id int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	rsp := idResponse{ID: id}
	err := json.NewEncoder(w).Encode(rsp)
	if err != nil {
		log.Println(err)
	}
}

type tasksResponse struct {
	Tasks []task.Task `json:"tasks"`
}

func sendTasks(w http.ResponseWriter, tasks []task.Task) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	rsp := tasksResponse{Tasks: tasks}
	err := json.NewEncoder(w).Encode(rsp)
	if err != nil {
		log.Println(err)
	}
}

func sendTask(w http.ResponseWriter, task task.Task) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Println(err)
	}
}

func sendOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{}")
}
