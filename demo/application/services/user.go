package services

import (
	"strings"
	"time"

	"github.com/duycs/demo-go/demo/entities"
	"github.com/duycs/demo-go/demo/infrastructure/helpers"
)

type UserReader interface {
	FindByID(id uint32) (*entities.User, error)
	Search(query string) ([]*entities.User, error)
	List() ([]*entities.User, error)
}

type UserWriter interface {
	Add(e *entities.User) (uint32, error)
	Update(e *entities.User) error
	Delete(id uint32) error
}

type UserRepository interface {
	UserReader
	UserWriter
}

type UserUseCase interface {
	GetUser(id uint32) (*entities.User, error)
	SearchUsers(query string) ([]*entities.User, error)
	ListUsers() ([]*entities.User, error)
	CreateUser(email, password, firstName, lastName string) (uint32, error)
	UpdateUser(e *entities.User) error
	DeleteUser(id uint32) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) CreateUser(email, password, firstName, lastName string) (uint32, error) {
	e, err := entities.CreateUser(email, password, firstName, lastName)
	if err != nil {
		return e.ID, err
	}
	return s.repo.Add(e)
}

func (s *UserService) GetUser(id uint32) (*entities.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) SearchUsers(query string) ([]*entities.User, error) {
	return s.repo.Search(strings.ToLower(query))
}

func (s *UserService) ListUsers() ([]*entities.User, error) {
	return s.repo.List()
}

func (s *UserService) DeleteUser(id uint32) error {
	u, err := s.GetUser(id)
	if u == nil {
		return helpers.ErrNotFound
	}
	if err != nil {
		return err
	}
	if len(u.Tasks) > 0 {
		return helpers.ErrCannotBeDeleted
	}
	return s.repo.Delete(id)
}

func (s *UserService) UpdateUser(e *entities.User) error {
	err := e.Validate()
	if err != nil {
		return helpers.ErrInvalidEntity
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
