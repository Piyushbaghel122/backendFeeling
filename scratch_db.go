package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"go-api/src/modules"
)

func main() {
	dsn := "root:password123@tcp(localhost:3306)/freelancing?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var users []modules.User
	db.Find(&users)
	fmt.Printf("Before: Found %d users\n", len(users))

	newUser := modules.User{
		UserName: "testuser",
		Email:    "test@test.com",
		Password: "password",
	}
	db.Create(&newUser)
	fmt.Printf("Created user with ID: %d\n", newUser.UserId)

	db.Find(&users)
	fmt.Printf("After: Found %d users\n", len(users))
}
