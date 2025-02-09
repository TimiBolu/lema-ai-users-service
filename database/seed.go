package database

import (
	"fmt"
	"reflect"
	"time"

	"github.com/TimiBolu/lema-ai-users-service/models"
	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already seeded")
		return
	}

	// Seed 50 users
	fkUUID := faker.UUID{}
	users := make([]models.User, 50)
	for i := range users {
		user := models.User{}
		uuidStr, _ := fkUUID.Hyphenated(reflect.Value{})
		user.ID = uuidStr.(string)
		user.FirstName = faker.FirstName()
		user.LastName = faker.LastName()
		user.Email = faker.Email()
		users[i] = user
	}

	// Seed addresses for each user
	addresses := make([]models.Address, 50)
	for i, user := range users {
		reakAddress := faker.GetRealAddress()
		uuidStr, _ := fkUUID.Hyphenated(reflect.Value{})

		address := models.Address{
			ID:      uuidStr.(string),
			UserID:  user.ID,
			Street:  reakAddress.Address,
			City:    reakAddress.City,
			State:   reakAddress.State,
			ZipCode: reakAddress.PostalCode,
		}
		addresses[i] = address
	}

	// Seed 4 posts per user
	posts := make([]models.Post, 0)
	fkLorem := faker.Lorem{}
	for _, user := range users {
		for i := 0; i < 4; i++ {
			uuidStr, _ := fkUUID.Hyphenated(reflect.Value{})
			lorem, _ := fkLorem.Paragraph(reflect.Value{})
			post := models.Post{
				ID:        uuidStr.(string),
				UserID:    user.ID,
				Title:     faker.Sentence(),
				Body:      lorem.(string),
				CreatedAt: time.Now(),
			}
			posts = append(posts, post)
		}
	}

	db.Create(&users)
	db.Create(&addresses)
	db.Create(&posts)

	fmt.Println("Database seeded successfully!")
}
