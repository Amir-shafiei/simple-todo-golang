package handlers

import (
	"awesomeProject19/models"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (u *UserHandler) Login(c *gin.Context) {
	var user models.User
	var existing models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := u.DB.Where("username = ?", user.Username).First(&existing)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "کاربر پیدا نشد"})
		return
	}
	er := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(user.Password))
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": existing.Username,
		"id":       existing.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
func (u *UserHandler) Register(c *gin.Context) {
	var user models.User
	var existing models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := u.DB.Where("username = ?", user.Username).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Already Exists"})
		return
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(password)
	u.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Signed Up Successfully"})

}
func (u *UserHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout success"})

}
