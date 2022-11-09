package repository

import (
	"errors"
	"time"

	"github.com/duycs/demo-go/demo/entities"
	"github.com/jinzhu/gorm"
)

type TaskContext struct {
	db *gorm.DB
}

func NewTaskContext(db *gorm.DB) *TaskContext {
	return &TaskContext{
		db: db,
	}
}

func (g *TaskContext) List(u *entities.Task) (*[]entities.Task, error) {
	var err error
	tasks := []entities.Task{}
	err = g.db.Debug().Model(&entities.Task{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]entities.Task{}, err
	}
	return &tasks, err
}

func (g *TaskContext) FindByID(u *entities.Task, uid uint32) (*entities.Task, error) {
	var err error
	err = g.db.Debug().Model(entities.Task{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entities.Task{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entities.Task{}, errors.New("Task Not Found")
	}
	return u, err
}

func (g *TaskContext) Search(u *entities.Task, query string) (*entities.Task, error) {
	var err error
	err = g.db.Debug().Model(entities.Task{}).Where("title = ?", query).Or("description = ?", query).Take(&u).Error
	if err != nil {
		return &entities.Task{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entities.Task{}, errors.New("Task Not Found")
	}
	return u, err
}

func (g *TaskContext) Add(u *entities.Task) (*entities.Task, error) {
	var err error
	err = g.db.Create(&u).Error
	if err != nil {
		return &entities.Task{}, err
	}
	return u, nil
}

func (g *TaskContext) Update(u *entities.Task, id uint32) (*entities.Task, error) {
	g.db = g.db.Debug().Model(&entities.Task{}).Where("id = ?", id).Take(&entities.Task{}).UpdateColumns(
		map[string]interface{}{
			"title":                u.Title,
			"description":          u.Description,
			"estimation_in_second": u.EstimationInSecond,
			"update_at":            time.Now(),
		},
	)
	if g.db.Error != nil {
		return &entities.Task{}, g.db.Error
	}

	err := g.db.Debug().Model(&entities.Task{}).Where("id = ?", id).Take(&u).Error
	if err != nil {
		return &entities.Task{}, err
	}
	return u, nil
}

func (g *TaskContext) Delete(u *entities.Task, uid uint32) (int64, error) {

	g.db = g.db.Debug().Model(&entities.Task{}).Where("id = ?", uid).Take(&entities.Task{}).Delete(&entities.Task{})

	if g.db.Error != nil {
		return 0, g.db.Error
	}
	return g.db.RowsAffected, nil
}
