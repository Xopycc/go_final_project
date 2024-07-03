package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"go-final-project/internal/task"
)

const tasksLimit = 10 // Переименованная константа

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) TaskAdd(t task.Task) (int, error) {
	res, err := r.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		return 0, fmt.Errorf("wrong query to db: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add error: %w", err)
	}
	return int(id), nil
}

func (r *Repository) TasksGet(t task.Task, search string) ([]task.Task, error) {

	tasks := []task.Task{}

	search = `%` + search + `%`

	rows, err := r.db.Query(`SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date ASC LIMIT :limit;`,
		sql.Named("search", search),
		sql.Named("limit", tasksLimit),
	)
	if err != nil {
		return []task.Task{}, fmt.Errorf("wrong query to db: %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("rows close error:", err)
		}
	}()

	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return tasks, fmt.Errorf("scan rows err: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return tasks, fmt.Errorf("rows err: %w", err)
	}

	return tasks, nil
}

func (r *Repository) TaskGet(id string) (task.Task, error) {
	row := r.db.QueryRow("SELECT * FROM scheduler WHERE id = :id;",
		sql.Named("id", id),
	)
	t := task.Task{}

	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, task.ErrNotFound
		}
		return t, fmt.Errorf("scan row err: %w", err)
	}

	return t, nil
}

func (r *Repository) TaskUpdate(tsk task.Task) error {
	res, err := r.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", tsk.ID),
		sql.Named("date", tsk.Date),
		sql.Named("title", tsk.Title),
		sql.Named("comment", tsk.Comment),
		sql.Named("repeat", tsk.Repeat),
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("failed to get rows affected: %w", task.ErrNotFound)
	}
	return nil
}

func (r *Repository) TaskDelete(tskID string) error {
	_, err := r.db.Exec("DELETE FROM scheduler where id = :id",
		sql.Named("id", tskID),
	)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
