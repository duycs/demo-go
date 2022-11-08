package repository

import (
	"errors"
	"time"

	"github.com/duycs/demo-go/demo/domain/entity"
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

func (g *TaskContext) List(u *entity.Task) (*[]entity.Task, error) {
	var err error
	tasks := []entity.Task{}
	err = g.db.Debug().Model(&entity.Task{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]entity.Task{}, err
	}
	return &tasks, err
}

func (g *TaskContext) FindByID(u *entity.Task, uid uint32) (*entity.Task, error) {
	var err error
	err = g.db.Debug().Model(entity.Task{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entity.Task{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entity.Task{}, errors.New("Task Not Found")
	}
	return u, err
}

func (g *TaskContext) Search(u *entity.Task, query string) (*entity.Task, error) {
	var err error
	err = g.db.Debug().Model(entity.Task{}).Where("title = ?", query).Or("description = ?", query).Take(&u).Error
	if err != nil {
		return &entity.Task{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entity.Task{}, errors.New("Task Not Found")
	}
	return u, err
}

func (g *TaskContext) Add(u *entity.Task) (*entity.Task, error) {
	var err error
	err = g.db.Create(&u).Error
	if err != nil {
		return &entity.Task{}, err
	}
	return u, nil
}

func (g *TaskContext) Update(u *entity.Task, id uint32) (*entity.Task, error) {
	g.db = g.db.Debug().Model(&entity.Task{}).Where("id = ?", id).Take(&entity.Task{}).UpdateColumns(
		map[string]interface{}{
			"title":                u.Title,
			"description":          u.Description,
			"estimation_in_second": u.EstimationInSecond,
			"update_at":            time.Now(),
		},
	)
	if g.db.Error != nil {
		return &entity.Task{}, g.db.Error
	}

	err := g.db.Debug().Model(&entity.Task{}).Where("id = ?", id).Take(&u).Error
	if err != nil {
		return &entity.Task{}, err
	}
	return u, nil
}

func (g *TaskContext) Delete(u *entity.Task, uid uint32) (int64, error) {

	g.db = g.db.Debug().Model(&entity.Task{}).Where("id = ?", uid).Take(&entity.Task{}).Delete(&entity.Task{})

	if g.db.Error != nil {
		return 0, g.db.Error
	}
	return g.db.RowsAffected, nil
}
