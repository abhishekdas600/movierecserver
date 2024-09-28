
package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`                     
	Name  string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`      
	Watchlist []Watchlist `json:"watchlist" db:"-"`   
	CreatedAt time.Time `json:"created_at" db:"created_at"` 
}
