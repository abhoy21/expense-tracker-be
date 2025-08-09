package main

import (
	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/abhoy21/expense-tracker/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitilazeDBConnection()
}

func main() {
	initializers.DBClient.AutoMigrate(&models.User{}, &models.Category{}, &models.Transaction{})
}
