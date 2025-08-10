package controllers

import (
	"net/http"
	"strconv"

	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/abhoy21/expense-tracker/models"
	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	var body struct {
		Amount      float64 `gorm:"not null" binding:"required"`
		Description string  `gorm:"type:varchar(255)"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount is required"})
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
	transaction := models.Transaction{
		Amount:      body.Amount,
		Description: body.Description,
		UserID:      user.ID,
		CategoryID:  &categoryIdUint,
	}

	if err := initializers.DBClient.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Transaction created successfully",
		"category": transaction,
	})
}

func UpdateTransaction(c *gin.Context) {
	var body struct {
		Amount      float64 `binding:"required"`
		Description string  `binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount is required"})
		return
	}

	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}
	user := u.(models.User)


	categoryIdStr := c.Param("categoryId")
	categoryIdUint64, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	categoryIdUint := uint(categoryIdUint64)


	transactionIdStr := c.Param("id")
	transactionIdUint64, err := strconv.ParseUint(transactionIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	transactionIdUint := uint(transactionIdUint64)
	transaction := models.Transaction{
		Amount:      body.Amount,
		Description: body.Description,
		UserID:      user.ID,
		CategoryID:  &categoryIdUint,
	}

	result := initializers.DBClient.Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ? AND id = ?", user.ID, categoryIdUint, transactionIdUint).
		Updates(&transaction)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating transaction"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully", "transaction": transaction})
}

func DeleteTransaction(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}
	user := u.(models.User)


	categoryIdStr := c.Param("categoryId")
	categoryIdUint64, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	categoryIdUint := uint(categoryIdUint64)


	transactionIdStr := c.Param("id")
	transactionIdUint64, err := strconv.ParseUint(transactionIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	transactionIdUint := uint(transactionIdUint64)

	result := initializers.DBClient.
		Where("user_id = ? AND category_id = ? AND id = ?", user.ID, categoryIdUint, transactionIdUint).
		Delete(&models.Transaction{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting transaction"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully", "rowsAffected": result.RowsAffected})
}
