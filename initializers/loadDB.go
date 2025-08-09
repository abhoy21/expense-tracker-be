package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBClient *gorm.DB

func InitilazeDBConnection() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DBClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error initiating Database Connection: %v", err)
	}
}
