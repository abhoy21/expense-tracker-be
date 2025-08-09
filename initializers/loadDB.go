package initializers

import (
	"log"
	"os"
	"time"

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

	// Get generic DB object for connection pool config
	sqlDB, err := DBClient.DB()
	if err != nil {
		log.Fatalf("Error getting DB from GORM: %v", err)
	}

	// Pool settings
	sqlDB.SetMaxOpenConns(25)                  // max number of open connections
	sqlDB.SetMaxIdleConns(25)                  // keep idle connections ready
	sqlDB.SetConnMaxLifetime(5 * time.Minute)  // recycle connections periodically
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)  // close idle connections after 1 min
}
