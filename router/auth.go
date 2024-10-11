package router

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/abhishekdas600/movierecserver/db"
	"github.com/abhishekdas600/movierecserver/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

func GetAuthCallbackFunction(r *gin.Engine) {
	r.GET("/auth/:provider", beginAuthHandler)
	r.GET("/auth/:provider/callback", HandleOAuthCallback)
	r.GET("/auth/logout", Logout)
}

func beginAuthHandler(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func HandleOAuthCallback(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("Failed to complete user authentication: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	if err := processOAuthUser(c, user); err != nil {
		log.Printf("Error processing OAuth user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	http.Redirect(c.Writer, c.Request, "https://movierec-client.vercel.app", http.StatusFound)
}

func processOAuthUser(c *gin.Context, user goth.User) error {
	
	existingUser, err := getUserByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var userID int
	if existingUser != nil {
		log.Printf("User already exists: %s", existingUser.Email)
		userID = existingUser.ID
	} else {
		
		newUser := models.User{
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: time.Now(),
		}

		
		if err := db.DB.Create(&newUser).Error; err != nil {
			return err
		}
		userID = newUser.ID
		log.Printf("New user created with ID: %d", userID)
	}

	
	session := sessions.Default(c)
	session.Set("email", user.Email)
	session.Set("name", user.Name)
	session.Set("user_id", userID)
	session.Set("logged_in", true)

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}

func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
