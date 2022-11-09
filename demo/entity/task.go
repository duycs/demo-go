package entity

import (
	"errors"
	"time"
)

type Task struct {
	ID                 uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Title              string    `gorm:"size:255;not null;unique" json:"title"`
	Description        string    `gorm:"size:255;not null;unique" json:"description"`
	EstimationInSecond int       `gorm:"not null;" json:"estimation_in_second"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsDelete           bool
}

func CreateTask(title, description string, estimationInSecond int) (*Task, error) {
	t := &Task{
		Title:              title,
		Description:        description,
		EstimationInSecond: estimationInSecond,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("Title can not empty")
	}

	return nil
}