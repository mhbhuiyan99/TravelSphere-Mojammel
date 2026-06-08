package models

import "time"

type WishlistItem struct {
	ID          string
	CountryName string
	Note        string
	Status      string // planned or visited
	CreatedAt   time.Time
}
