package models

import "time"

type Watchlist struct {
	ID        int       `json:"id" gorm:"primaryKey"` 
	UserID    int       `json:"user_id" gorm:"not null"`
	TMDBID    int       `json:"tmdb_id" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"` 
}

type Favourites struct {
	ID        int       `json:"id" gorm:"primaryKey"` 
	UserID    int       `json:"user_id" gorm:"not null"`
	TMDBID    int       `json:"tmdb_id" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"` 
}