package main

import (
	"ToDo/internal/notes"
	"ToDo/internal/user"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Включаем логирование запросов в консоль.
	})
	if err != nil {
		panic(err)
	}

	db.Migrator().DropTable(&notes.Note{}, &user.User{})
	err = db.AutoMigrate(&user.User{}, &notes.Note{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations complete")
}
