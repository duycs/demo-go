package repository

import (
	"errors"
	"log"
	"time"

	"github.com/duycs/demo-go/demo/domain/entity"
	"github.com/duycs/demo-go/demo/infrastructure/helpers"
	"github.com/jinzhu/gorm"
)

type UserContext struct {
	db *gorm.DB
}

func NewUserContext(db *gorm.DB) *UserContext {
	return &UserContext{
		db: db,
	}
}

func (g *UserContext) List(u *entity.User) (*[]entity.User, error) {
	var err error
	users := []entity.User{}
	err = g.db.Debug().Model(&entity.User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]entity.User{}, err
	}
	return &users, err
}

func (g *UserContext) FindByID(u *entity.User, uid uint32) (*entity.User, error) {
	var err error
	err = g.db.Debug().Model(entity.User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entity.User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entity.User{}, errors.New("User Not Found")
	}
	return u, err
}

func (g *UserContext) Search(u *entity.User, query string) (*entity.User, error) {
	var err error
	err = g.db.Debug().Model(entity.User{}).Where("email = ?", query).Or("first_name = ?", query).Or("last_name = ?", query).Take(&u).Error
	if err != nil {
		return &entity.User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entity.User{}, errors.New("User Not Found")
	}
	return u, err
}

func (g *UserContext) Add(u *entity.User) (*entity.User, error) {
	var err error
	err = g.db.Create(&u).Error
	if err != nil {
		return &entity.User{}, err
	}
	return u, nil
}

func (g *UserContext) Update(u *entity.User, uid uint32) (*entity.User, error) {
	err := BeforeSave(u)
	if err != nil {
		log.Fatal(err)
	}
	g.db = g.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&entity.User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"email":      u.Email,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"update_at":  time.Now(),
		},
	)
	if g.db.Error != nil {
		return &entity.User{}, g.db.Error
	}

	err = g.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entity.User{}, err
	}
	return u, nil
}

func (g *UserContext) Delete(u *entity.User, uid uint32) (int64, error) {

	g.db = g.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&entity.User{}).Delete(&entity.User{})

	if g.db.Error != nil {
		return 0, g.db.Error
	}
	return g.db.RowsAffected, nil
}

func BeforeSave(u *entity.User) error {
	hashedPassword, err := helpers.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
