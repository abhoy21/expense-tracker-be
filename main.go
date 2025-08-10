package main

import (
	"os"

	"github.com/abhoy21/expense-tracker/controllers"
	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/abhoy21/expense-tracker/middleware"
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

	// User Routes
	router.POST("/sign-up", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/logout", controllers.Logout)
	router.GET("/user/details", middleware.ValidateAuth, controllers.GetUserDetails)

	// Category Routes
	router.POST("/create/category", middleware.ValidateAuth, controllers.CreateCategory)
	router.PATCH("/update/category/:id", middleware.ValidateAuth, controllers.UpdateCategory)
	router.DELETE("/delete/category/:id", middleware.ValidateAuth, controllers.DeleteCategory)

	// Transaction Routes
	router.POST("/create/transaction/:id", middleware.ValidateAuth, controllers.CreateTransaction)
	router.PATCH("/update/transaction/:categoryId/:id", middleware.ValidateAuth, controllers.UpdateTransaction)
	router.DELETE("/delete/transaction/:categoryId/:id", middleware.ValidateAuth, controllers.DeleteTransaction)

	router.Run(":" + port)
}
