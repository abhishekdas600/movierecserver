package user

import (
	"net/http"

	"github.com/abhishekdas600/movierecserver/db"
	"github.com/abhishekdas600/movierecserver/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func GetUserDetails(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email")

	if email == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}


	emailStr := email.(string)


	var user models.User
	if err := db.GetDB().Where("email = ?", emailStr).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details", "details": err.Error()})
		}
		return
	}


	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, user)
}

