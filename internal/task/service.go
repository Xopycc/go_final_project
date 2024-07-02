package task

import (
	"fmt"
	"time"
)

type Repository interface {
	TaskAdd(t Task) (int, error)
	TasksGet(t Task, search string) ([]Task, error)
	TaskGet(id string) (Task, error)
	TaskUpdate(t Task) error
	TaskDelete(id string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	u := Service{repo: repo}
	return &u
}

func (s *Service) Create(tsk Task) (int, error) {

	err := validateData(&tsk)
	if err != nil {
		return 0, err
	}

	if tsk.Repeat == "" || tsk.Repeat == " " {
		tsk.Date = time.Now().Format(DateFormat)
		taskID, err := s.repo.TaskAdd(tsk)
		if err != nil {
			return 0, fmt.Errorf("task - Create: %w", err)

		}
		return taskID, nil
	}
	now := time.Now().Truncate(24 * time.Hour)

	nowText := now.Format(DateFormat)
	if tsk.Date < nowText {
		tsk.Date, err = NextDate(now, tsk.Date, tsk.Repeat)
		if err != nil {
			return 0, fmt.Errorf("%w: repeat is empty", ErrFormat)
		}
	}

	taskID, err := s.repo.TaskAdd(tsk)
	if err != nil {
		return 0, fmt.Errorf("task - Create: %w", err)
	}
	return taskID, nil
}

func (s *Service) GetTasks(tsk Task, search string) ([]Task, error) {
	res, err := s.repo.TasksGet(tsk, search)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetTask(id string) (Task, error) {
	tsk, err := s.repo.TaskGet(id)
	if err != nil {
		return tsk, err
	}
	return tsk, nil
}

func (s *Service) UpdateTask(tsk Task) error {

	if tsk.ID == "" || tsk.ID == " " {
		return fmt.Errorf("%w: ID is empty", ErrFormat)
	}

	err := validateData(&tsk)
	if err != nil {
		return err
	}

	if tsk.Repeat == "" || tsk.Repeat == " " {
		tsk.Date = time.Now().Format(DateFormat)
		err := s.repo.TaskUpdate(tsk)
		if err != nil {
			return fmt.Errorf("task - Update: %w", err)
		}
		return nil
	}
	now := time.Now().Truncate(24 * time.Hour)

	nowText := now.Format(DateFormat)
	if tsk.Date < nowText {
		tsk.Date, err = NextDate(now, tsk.Date, tsk.Repeat)
		if err != nil {
			return fmt.Errorf("%w: repeat is empty", ErrFormat)
		}
	}

	err = s.repo.TaskUpdate(tsk)
	if err != nil {
		return fmt.Errorf("%w: ", ErrFormat)
	}
	return nil
}

func (s *Service) TaskDone(tskID string) error {
	tsk, err := s.repo.TaskGet(tskID)
	if err != nil {
		return fmt.Errorf("cant get task: %w", err)
	}
	if tsk.Repeat == "" {
		err = s.repo.TaskDelete(tskID)
		if err != nil {
			return fmt.Errorf("cant delete task: %w", err)
		}
	} else {
		now := time.Now().Truncate(24 * time.Hour)
		tsk.Date, err = NextDate(now, tsk.Date, tsk.Repeat)
		if err != nil {
			return fmt.Errorf("failed to get next date: %w", ErrFormat)
		}
		err = s.repo.TaskUpdate(tsk)
		if err != nil {
			return fmt.Errorf("failed to update task: %w", ErrFormat)
		}
	}
	return nil
}

func (s *Service) TaskDelete(tskID string) error {
	_, err := s.repo.TaskGet(tskID)
	if err != nil {
		return fmt.Errorf("cant get task: %w", err)
	}

	err = s.repo.TaskDelete(tskID)
	if err != nil {
		return fmt.Errorf("cant delete task: %w", err)
	}

	return nil
}
