package database

import (
	"fmt"
	"log"
	"os"

	"github.com/pratikjagrut/myreads-backend/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_pass, db_host, db_port, db_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database! |n", err)
	}

	log.Println("Connected to the database successfully")
	log.Println("Running Migrations!")

	// Add migrations
	err = db.AutoMigrate(&models.User{}, &models.Book{})
	if err != nil {
		log.Fatal("Failed to add migrations")
	}
	Database = DbInstance{Db: db}
}
