package seed

import (
	"log"

	"github.com/duycs/demo-go/demo/entities"
	"github.com/jinzhu/gorm"
)

var users = []entities.User{
	entities.User{
		FirstName: "Duy",
		Email:     "duycs@gmail.com",
		Password:  "password",
	},
	entities.User{
		FirstName: "Phuong",
		Email:     "phuong@gmail.com",
		Password:  "password",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&entities.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&entities.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&entities.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
