package main

import (
	"os"

	"github.com/abhoy21/expense-tracker/controllers"
	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitilazeDBConnection()
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

	// User SignUp, Login, Logout Routes
	router.POST("/sign-up", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/logout", controllers.Logout)

	router.Run(":" + port)
}
