package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/duycs/demo-go/demo/models"
)

var users = []models.User{
	models.User{
		Name: "Duy",
		Email:    "duycs@gmail.com",
		Password: "password",
	},
	models.User{
		Name: "Phuong",
		Email:    "phuong@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
