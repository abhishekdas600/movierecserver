package router

import (
	"context"
	"log"
	"github.com/abhishekdas600/movierecserver/models"
	"net/http"
	"time"

	"database/sql"
	"github.com/abhishekdas600/movierecserver/db"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func GetAuthCallbackFunction(r *gin.Engine) {
	r.GET("/auth/:provider", beginAuthHandler)
	r.GET("/auth/:provider/callback", HandleOAuthCallback)
	r.GET("auth/logout", Logout)
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
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
    
	if err := processOAuthUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	http.Redirect(c.Writer, c.Request, "http://localhost:8080", http.StatusFound)
}

func processOAuthUser(c *gin.Context, user goth.User) error {
	var dbConn *sql.DB
	if err := db.InitializePostgreSQL(&dbConn); err != nil {
		return err
	}
	defer dbConn.Close()

	existingUser, err := getUserByEmail(dbConn, user.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		log.Printf("User already exists: %s", existingUser.Email)
	} else {
		newUser := models.User{
			Email:     user.Email,
			Name: user.Name,
			CreatedAt: time.Now(),
		}


		query := `INSERT INTO users (email, name, created_at) 
				  VALUES ($1, $2, $3) RETURNING id`
		var userID int
		err = dbConn.QueryRow(query, newUser.Email, newUser.Name, newUser.CreatedAt).Scan(&userID)
		if err != nil {
			return err
		}

		
	}

	session := sessions.Default(c) 
	session.Set("email", user.Email)
	session.Set("name", user.Name)
	session.Set("logged_in", true) 

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}


func getUserByEmail(dbConn *sql.DB, email string) (*models.User, error) {
    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    stmt, err := dbConn.PrepareContext(ctx, "SELECT email, name, created_at FROM users WHERE email = $1")
    if err != nil {
        return nil, err
    }
    defer stmt.Close() 

    err = stmt.QueryRowContext(ctx, email).Scan(&user.Email, &user.Name, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil 
        }
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