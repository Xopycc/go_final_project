package task

import (
	"errors"
	"fmt"
	"time"
)

const DateFormat = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var (
	ErrFormat   = errors.New("task format error")
	ErrNotFound = errors.New("not found")
)

func validateData(tsk *Task) error {

	if tsk.Title == "" || tsk.Title == " " {
		return fmt.Errorf("%w: title is empty", ErrFormat)
	}

	if tsk.Date == "" {
		tsk.Date = time.Now().Format(DateFormat)
	}

	_, err := time.Parse(DateFormat, tsk.Date)
	if err != nil {
		return fmt.Errorf("%w: wrong data format", ErrFormat)
	}

	return nil

}
