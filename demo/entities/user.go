package entities

import (
	"time"

	"github.com/duycs/demo-go/demo/infrastructure/helpers"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Account   string    `gorm:"size:255;not null;unique" json:"account"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	FirstName string    `gorm:"size:255;not null;unique" json:"first_name"`
	LastName  string    `gorm:"size:255;not null;unique" json:"last_name"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsDelete  int       `gorm:"-"`
	Tasks     []uint32
}

func CreateUser(email, password, firstName, lastName string) (*User, error) {
	u := &User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
	}
	pwd, err := helpers.GeneratePassword(password)
	if err != nil {
		return nil, err
	}
	u.Password = pwd
	err = u.Validate()
	if err != nil {
		return nil, helpers.ErrInvalidEntity
	}
	return u, nil
}

func (u *User) AddTask(id uint32) error {
	_, err := u.GetTask(id)
	if err == nil {
		return helpers.ErrNotFound
	}
	u.Tasks = append(u.Tasks, id)
	return nil
}

func (u *User) RemoveTask(id uint32) error {
	for i, j := range u.Tasks {
		if j == id {
			u.Tasks = append(u.Tasks[:i], u.Tasks[i+1:]...)
			return nil
		}
	}
	return helpers.ErrNotFound
}

func (u *User) GetTask(id uint32) (uint32, error) {
	for _, v := range u.Tasks {
		if v == id {
			return id, nil
		}
	}
	return id, helpers.ErrNotFound
}

func (u *User) Validate() error {
	if u.Password == "" {
		return helpers.FormatError("Required Password")
	}
	if u.Email == "" {
		return helpers.FormatError("Required Email")
	}
	if err := helpers.ValidateEmail(u.Email); err != nil {
		return helpers.FormatError("Invalid Email")
	}
	return nil
}
