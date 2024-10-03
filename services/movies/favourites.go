package movies

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abhishekdas600/movierecserver/db"
	"github.com/abhishekdas600/movierecserver/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserFavourites(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	ctx := context.Background()

	var favourites []models.Favourites
	err := db.GetDB().WithContext(ctx).Where("user_id = ?", userID.(int)).Find(&favourites).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favourites"})
		return
	}

	c.JSON(http.StatusOK, favourites)
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

	ctx := context.Background()

	favouritesItem := models.Favourites{
		UserID:    userID.(int),
		TMDBID:    tmdbID,
		CreatedAt: time.Now(),
	}

	err = db.GetDB().WithContext(ctx).Create(&favouritesItem).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add movie to favourites: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie added to favourites successfully"})
}

func RemoveMovieFromFavourites(c *gin.Context) {
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

	ctx := context.Background()

	var favourite models.Favourites
	err = db.GetDB().WithContext(ctx).Where("user_id = ? AND tmdb_id = ?", userID.(int), tmdbID).Delete(&favourite).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in favourites"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove movie from favourites: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie removed from favourites successfully"})
}
