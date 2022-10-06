package seed

import (
	"log"

	"github.com/duycs/demo-go/demo/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Account:  "Duy",
		Email:    "duycs@gmail.com",
		Password: "password",
	},
	models.User{
		Account:  "Phuong",
		Email:    "phuong@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
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
