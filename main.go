package main

import (
	"os"

	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	router := gin.Default()
	port := os.Getenv(("PORT"))

	if port == "" {
		port = "8080"
	}

	router.GET("/", func(c *gin.Context) {
    c.String(200, "Bankai")
	})

	router.Run(":" + port)
}
