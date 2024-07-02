package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-final-project/internal/task"
)

type Api struct {
	task *task.Service
}

func New(tsk *task.Service) *Api {
	return &Api{task: tsk}
}

func (a *Api) TaskCreate(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tsk task.Task
	if err := json.Unmarshal(buf.Bytes(), &tsk); err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskID, err := a.task.Create(tsk)
	if err != nil {
		if errors.Is(err, task.ErrFormat) {
			sendErr(w, err.Error(), http.StatusBadRequest)
		} else {
			sendErr(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	sendID(w, taskID)

}

func (a *Api) NextDate(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	now, err := time.Parse(task.DateFormat, nowStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")

	repeat := r.FormValue("repeat")

	nextDate, err := task.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}

}

func (a *Api) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tsk task.Task
	search := r.FormValue("search")
	tasks, err := a.task.GetTasks(tsk, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	sendTasks(w, tasks)
}

func (a *Api) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendErr(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		sendErr(w, "Идентификатор должен быть числом", http.StatusBadRequest)
		return
	}
	tsk, err := a.task.GetTask(id)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			sendErr(w, "Задача не найдена", http.StatusNotFound)
			return
		}
		sendErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendTask(w, tsk)
}

func (a *Api) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tsk task.Task
	if err := json.Unmarshal(buf.Bytes(), &tsk); err != nil {
		sendErr(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.task.UpdateTask(tsk)
	if err != nil {
		if errors.Is(err, task.ErrFormat) {
			sendErr(w, err.Error(), http.StatusBadRequest)
		} else {
			sendErr(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	sendOK(w)
}

func (a *Api) TaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendErr(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	err := a.task.TaskDone(id)
	if err != nil {
		sendErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendOK(w)
}

func (a *Api) TaskDelete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendErr(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	err := a.task.TaskDelete(id)
	if err != nil {
		sendErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendOK(w)
}