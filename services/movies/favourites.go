package movies

import (
	"database/sql"
	"fmt"
	"github.com/abhishekdas600/movierecserver/db"
	"github.com/abhishekdas600/movierecserver/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserFavourites(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	var dbConn *sql.DB
	if err := db.InitializePostgreSQL(&dbConn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer dbConn.Close()

	favourites, err := getFavouritesByUserID(dbConn, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favourites"})
		return
	}

	c.JSON(http.StatusOK, favourites)
}

func getFavouritesByUserID(db *sql.DB, userID int) ([]models.Favourites, error) {
	query := `SELECT id, user_id, tmdb_id, created_at FROM favourites WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favourites []models.Favourites
	for rows.Next() {
		var fv models.Favourites
		if err := rows.Scan(&fv.ID, &fv.UserID, &fv.TMDBID, &fv.CreatedAt); err != nil {
			return nil, err
		}
		favourites = append(favourites, fv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favourites, nil
}

func AddMovieToFavourites(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id") 
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	tmdbIDStr := c.Param("tmdb_id") 
	if tmdbIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TMDB ID is required"})
		return
	}

	tmdbID, err := strconv.Atoi(tmdbIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TMDB ID format"})
		return
	}

	db := c.MustGet("db").(*sql.DB) 

	favouritesItem := models.Favourites{
		UserID:   userID.(int),
		TMDBID:   tmdbID, 
		CreatedAt: time.Now(),
	}

	err = addMovieToFavourites(db, favouritesItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add movie to favourites: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie added to favourites successfully"})
}
func addMovieToFavourites(db *sql.DB, fv models.Favourites) error {
	query := `INSERT INTO favourites (user_id, tmdb_id, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, fv.UserID, fv.TMDBID, fv.CreatedAt)
	return err
}