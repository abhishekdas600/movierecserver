package user

import (
	"database/sql"
	
	"github.com/abhishekdas600/movierecserver/db"
	"github.com/abhishekdas600/movierecserver/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)


func GetUserDetails(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email") 

	if email == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	var dbConn *sql.DB
	if err := db.InitializePostgreSQL(&dbConn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer dbConn.Close()

	user, err := getUserByEmail(dbConn, email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	session.Set("user_id", user.ID) 
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getUserByEmail(db *sql.DB, email string) (*models.User, error) {
	query := `SELECT id, email, name, created_at FROM users WHERE email = $1`
	var user models.User
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil 
	} else if err != nil {
		return nil, err 
	}
	return &user, nil
}