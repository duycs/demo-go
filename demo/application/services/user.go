package services

import (
	"strings"
	"time"

	"github.com/duycs/demo-go/demo/domain/entity"
	"github.com/duycs/demo-go/demo/infrastructure/helpers"
)

type UserReader interface {
	FindByID(id uint32) (*entity.User, error)
	Search(query string) ([]*entity.User, error)
	List() ([]*entity.User, error)
}

type UserWriter interface {
	Add(e *entity.User) (uint32, error)
	Update(e *entity.User) error
	Delete(id uint32) error
}

type UserRepository interface {
	UserReader
	UserWriter
}

type UserUseCase interface {
	GetUser(id uint32) (*entity.User, error)
	SearchUsers(query string) ([]*entity.User, error)
	ListUsers() ([]*entity.User, error)
	CreateUser(email, password, firstName, lastName string) (uint32, error)
	UpdateUser(e *entity.User) error
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
	e, err := entity.CreateUser(email, password, firstName, lastName)
	if err != nil {
		return e.ID, err
	}
	return s.repo.Add(e)
}

func (s *UserService) GetUser(id uint32) (*entity.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) SearchUsers(query string) ([]*entity.User, error) {
	return s.repo.Search(strings.ToLower(query))
}

func (s *UserService) ListUsers() ([]*entity.User, error) {
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

func (s *UserService) UpdateUser(e *entity.User) error {
	err := e.Validate()
	if err != nil {
		return helpers.ErrInvalidEntity
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
