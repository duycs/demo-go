package services

import (
	"errors"

	"github.com/duycs/demo-go/demo/entities"
)

type AssignmentUseCase interface {
	Assign(u *entities.User, t *entities.Task) error
	Checkout(t *entities.Task) error
}

type AssignmentService struct {
	userService UserUseCase
	taskService TaskUseCase
}

func NewAssignmentService(u UserUseCase, t TaskUseCase) *AssignmentService {
	return &AssignmentService{
		userService: u,
		taskService: t,
	}
}

func (s *AssignmentService) Assign(u *entities.User, t *entities.Task) error {
	u, err := s.userService.GetUser(u.ID)
	if err != nil {
		return err
	}

	t, err = s.taskService.GetTask(t.ID)
	if err != nil {
		return err
	}

	if t.EstimationInSecond <= 0 {
		return errors.New("Estimation time expired")
	}

	err = u.AddTask(t.ID)
	if err != nil {
		return err
	}

	err = s.userService.UpdateUser(u)
	if err != nil {
		return err
	}

	t.Assigned()
	err = s.taskService.UpdateTask(t)
	if err != nil {
		return err
	}
	return nil
}

// checkout task was finished then remove this
func (s *AssignmentService) Checkout(t *entities.Task) error {
	t, err := s.taskService.GetTask(t.ID)
	if err != nil {
		return err
	}

	users, err := s.userService.ListUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		_, err := u.GetTask(t.ID)
		if err != nil {
			continue
		}

		if t.Status == entities.Done {
			err = u.RemoveTask(t.ID)
			if err != nil {
				return err
			}

			err = s.userService.UpdateUser(u)
			if err != nil {
				return err
			}
		}

		break
	}

	return nil
}
