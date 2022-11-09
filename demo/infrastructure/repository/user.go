package repository

import (
	"errors"
	"log"
	"time"

	"github.com/duycs/demo-go/demo/entities"
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

func (g *UserContext) List(u *entities.User) (*[]entities.User, error) {
	var err error
	users := []entities.User{}
	err = g.db.Debug().Model(&entities.User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]entities.User{}, err
	}
	return &users, err
}

func (g *UserContext) FindByID(u *entities.User, uid uint32) (*entities.User, error) {
	var err error
	err = g.db.Debug().Model(entities.User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entities.User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entities.User{}, errors.New("User Not Found")
	}
	return u, err
}

func (g *UserContext) Search(u *entities.User, query string) (*entities.User, error) {
	var err error
	err = g.db.Debug().Model(entities.User{}).Where("email = ?", query).Or("first_name = ?", query).Or("last_name = ?", query).Take(&u).Error
	if err != nil {
		return &entities.User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entities.User{}, errors.New("User Not Found")
	}
	return u, err
}

func (g *UserContext) Add(u *entities.User) (*entities.User, error) {
	var err error
	err = g.db.Create(&u).Error
	if err != nil {
		return &entities.User{}, err
	}
	return u, nil
}

func (g *UserContext) Update(u *entities.User, uid uint32) (*entities.User, error) {
	err := BeforeSave(u)
	if err != nil {
		log.Fatal(err)
	}
	g.db = g.db.Debug().Model(&entities.User{}).Where("id = ?", uid).Take(&entities.User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"email":      u.Email,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"update_at":  time.Now(),
		},
	)
	if g.db.Error != nil {
		return &entities.User{}, g.db.Error
	}

	err = g.db.Debug().Model(&entities.User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entities.User{}, err
	}
	return u, nil
}

func (g *UserContext) Delete(u *entities.User, uid uint32) (int64, error) {

	g.db = g.db.Debug().Model(&entities.User{}).Where("id = ?", uid).Take(&entities.User{}).Delete(&entities.User{})

	if g.db.Error != nil {
		return 0, g.db.Error
	}
	return g.db.RowsAffected, nil
}

func BeforeSave(u *entities.User) error {
	hashedPassword, err := helpers.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
