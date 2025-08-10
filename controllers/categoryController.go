package controllers

import (
	"net/http"
	"strconv"

	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/abhoy21/expense-tracker/models"
	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var body struct {
		Name string `binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	u, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Your are not authorized to perform this action"})
	}

	user := u.(models.User)
	category := models.Category{
    UserID: user.ID,
    Name:   body.Name,
	}

	if err := initializers.DBClient.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category created successfully",
		"category": category,
	})
}

func UpdateCategory(c *gin.Context) {
	var body struct {
		Name string `binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	u, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Your are not authorized to perform this action"})
	}

	user := u.(models.User)
	categoryIdStr := c.Param("id")
	categoryIdUint64, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
	}
	categoryIdUint := uint(categoryIdUint64)
	category := models.Category{
    UserID: user.ID,
    Name:   body.Name,
	}

	if err := initializers.DBClient.Model(&category).
	Where("user_id = ? AND id = ?", user.ID, categoryIdUint).
	Update("name", body.Name).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error updating category, please try again!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updates successfully!", "category": category})
}


func DeleteCategory(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	user := u.(models.User)

	categoryIdStr := c.Param("id")
	categoryIdUint64, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	categoryIdUint := uint(categoryIdUint64)

	result := initializers.DBClient.
		Where("user_id = ? AND id = ?", user.ID, categoryIdUint).
		Delete(&models.Category{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting category"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully", "rowsAffected": result.RowsAffected})
}
