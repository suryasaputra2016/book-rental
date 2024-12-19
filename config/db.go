package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/suryasaputra2016/book-rental/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connects to postgres database using gorm and returns db connection
func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// opean database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// migrate database
	migrateDatabase(db)

	log.Println("Database connected successfully.")
	return db
}

// migrateDatabase migrate database
func migrateDatabase(db *gorm.DB) {
	fmt.Println("Running migration")
	err := db.AutoMigrate(
		entity.User{},
		entity.Rent{},
		entity.Book{},
		entity.BookCopy{},
		entity.RentalHistory{},
	)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	Session := db.Session(&gorm.Session{PrepareStmt: true})
	if Session != nil {
		fmt.Println("Migration successful")
	}
}
