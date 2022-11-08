package seed

import (
	"log"

	"github.com/duycs/demo-go/demo/domain/entity"
	"github.com/jinzhu/gorm"
)

var users = []entity.User{
	entity.User{
		FirstName: "Duy",
		Email:     "duycs@gmail.com",
		Password:  "password",
	},
	entity.User{
		FirstName: "Phuong",
		Email:     "phuong@gmail.com",
		Password:  "password",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&entity.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&entity.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&entity.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
