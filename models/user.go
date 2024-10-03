
package models

import "time"

type User struct {
	ID         int         `json:"id" gorm:"primaryKey"` 
	Name       string      `json:"name" gorm:"not null"`
	Email      string      `json:"email" gorm:"unique;not null"` 
	Watchlist  []Watchlist `json:"watchlist" gorm:"foreignKey:UserID"` 
	Favourites []Favourites `json:"favourites" gorm:"foreignKey:UserID"` 
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
}