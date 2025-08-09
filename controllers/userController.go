package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/abhoy21/expense-tracker/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)



func SignUp(c *gin.Context){
	var body struct {
		Email    string `gorm:"not null;uniqueIndex"`
		Name 		string `gorm:"not null"`
    Password string `gorm:"not null"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if body.Email == "" || body.Password == "" || body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email, Name and password are required",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error hashing Password"})
	}

	user := models.User{Email: body.Email, Name: body.Name, Password: string(hash)}

	response := initializers.DBClient.Create(&user)

	if response.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Signing Up User"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User SignedUp Successfully!", "user": gin.H{"id": user.ID, "email": user.Email, "Name": user.Name}})
}


func Login(c *gin.Context) {
	var body struct {
		Email    string `gorm:"not null;uniqueIndex"`
    Password string `gorm:"not null"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
	}

	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
 	}
 	var user models.User
 	// initializers.DBClient.First(&user, "Email = ?", body.Email)

	//Optimizations -> Fetch only required fields (faster DB query)

	start := time.Now()
	if err := initializers.DBClient.
		Select("id", "email", "password").
		Where("email = ?", body.Email).
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	queryTime := time.Since(start)
	log.Printf("DB query for Login took: %v", queryTime)

 	if user.ID == 0 {
		c.JSON(404, gin.H{"error": "User Not Found, Invalid Credentials"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Password, Please try again with Correct Password"})
		return;
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(400, gin.H{"error": "Error signing Token!"})
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 7, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User Successfully Logged In!"})
}


func Logout(c *gin.Context){
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User successfully logged out"})
}
