package models

import "time"

type Watchlist struct {
	ID        int       `json:"id" db:"id"`                 
	UserID    int       `json:"user_id" db:"user_id"`      
	TMDBID    int       `json:"tmdb_id" db:"tmdb_id"`       
	CreatedAt time.Time `json:"created_at" db:"created_at"` 

}

type Favourites struct {
	ID        int       `json:"id" db:"id"`                 
	UserID    int       `json:"user_id" db:"user_id"`      
	TMDBID    int       `json:"tmdb_id" db:"tmdb_id"`       
	CreatedAt time.Time `json:"created_at" db:"created_at"` 

}