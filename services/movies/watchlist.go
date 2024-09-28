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

func GetUserWatchlist(c *gin.Context) {
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

	watchlist, err := getWatchlistByUserID(dbConn, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch watchlist"})
		return
	}

	c.JSON(http.StatusOK, watchlist)
}

func getWatchlistByUserID(db *sql.DB, userID int) ([]models.Watchlist, error) {
	query := `SELECT id, user_id, tmdb_id, created_at FROM watchlist WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var watchlist []models.Watchlist
	for rows.Next() {
		var wl models.Watchlist
		if err := rows.Scan(&wl.ID, &wl.UserID, &wl.TMDBID, &wl.CreatedAt); err != nil {
			return nil, err
		}
		watchlist = append(watchlist, wl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return watchlist, nil
}

func AddMovieToWatchlist(c *gin.Context) {
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

	watchlistItem := models.Watchlist{
		UserID:   userID.(int),
		TMDBID:   tmdbID, 
		CreatedAt: time.Now(),
	}

	err = addMovieToWatchlist(db, watchlistItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add movie to watchlist: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie added to watchlist successfully"})
}
func addMovieToWatchlist(db *sql.DB, wl models.Watchlist) error {
	query := `INSERT INTO watchlist (user_id, tmdb_id, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, wl.UserID, wl.TMDBID, wl.CreatedAt)
	return err
}