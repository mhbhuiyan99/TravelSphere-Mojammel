package models

import "time"

type WishlistItem struct {
	ID          string
	CountryName string
	Note        string
	Status      string // "Planned" or "Visited"
	CreatedAt   time.Time
}

var AllowedStatuses = map[string]bool{
	"Planned": true,
	"Visited": true,
}