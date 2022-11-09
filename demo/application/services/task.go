package services

import (
	"strings"
	"time"

	"github.com/duycs/demo-go/demo/entities"
	"github.com/duycs/demo-go/demo/infrastructure/helpers"
)

type TaskReader interface {
	FindByID(id uint32) (*entities.Task, error)
	Search(query string) ([]*entities.Task, error)
	List() ([]*entities.Task, error)
}

type TaskWriter interface {
	Add(e *entities.Task) (uint32, error)
	Update(e *entities.Task) error
	Delete(id uint32) error
}

type TaskRepository interface {
	TaskReader
	TaskWriter
}

type TaskUseCase interface {
	GetTask(id uint32) (*entities.Task, error)
	SearchTasks(query string) ([]*entities.Task, error)
	ListTasks() ([]*entities.Task, error)
	CreateTask(title string, description string, estimationInSecond int) (uint32, error)
	UpdateTask(e *entities.Task) error
	DeleteTask(id uint32) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(r TaskRepository) *TaskService {
	return &TaskService{
		repo: r,
	}
}

func (s *TaskService) CreateTask(title string, description string, estimationInSecond int) (uint32, error) {
	t, err := entities.CreateTask(title, description, estimationInSecond)
	if err != nil {
		return t.ID, err
	}
	return s.repo.Add(t)
}

func (s *TaskService) GetTask(id uint32) (*entities.Task, error) {
	t, err := s.repo.FindByID(id)
	if t == nil {
		return nil, helpers.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TaskService) SearchTasks(query string) ([]*entities.Task, error) {
	tasks, err := s.repo.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, helpers.ErrNotFound
	}
	return tasks, nil
}

func (s *TaskService) ListTasks() ([]*entities.Task, error) {
	tasks, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, helpers.ErrNotFound
	}
	return tasks, nil
}

func (s *TaskService) DeleteTask(id uint32) error {
	_, err := s.GetTask(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *TaskService) UpdateTask(e *entities.Task) error {
	err := e.Validate()
	if err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
