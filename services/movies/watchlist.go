package movies

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abhishekdas600/movierecserver/db"
	"github.com/abhishekdas600/movierecserver/models"
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

	var watchlist []models.Watchlist
	if err := db.GetDB().Where("user_id = ?", userID.(int)).Find(&watchlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch watchlist"})
		return
	}

	c.JSON(http.StatusOK, watchlist)
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

	watchlistItem := models.Watchlist{
		UserID:    userID.(int),
		TMDBID:    tmdbID,
		CreatedAt: time.Now(),
	}

	if err := db.GetDB().Create(&watchlistItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add movie to watchlist: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie added to watchlist successfully"})
}

func RemoveMovieFromWatchlist(c *gin.Context) {
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

	if err := db.GetDB().Where("user_id = ? AND tmdb_id = ?", userID.(int), tmdbID).Delete(&models.Watchlist{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove movie from watchlist: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie removed from watchlist successfully"})
}
