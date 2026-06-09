package services

import (
	"TravelSphere-Mojammel/models"
	"errors"
	"fmt"
	"sync"
	"time"
)

// store holds wishlists keyed by username - resets on server restart
var (
	store = map[string][]models.WishlistItem{}
	storeMu sync.RWMutex
)

// GetWishlist returns all wishlist items for a user
func GetWishlist(username string) []models.WishlistItem {
	storeMu.RLock()
	defer storeMu.Unlock()
	return store[username]

}

// AddToWishlist creates a new wishlist entry for a user
func AddToWishlist(username, countryName string) (models.WishlistItem, error) {
	if username == "" {
		return models.WishlistItem{}, errors.New("username required")
	}
	if countryName == "" {
		return models.WishlistItem{}, errors.New("country name required")
	}

	items := models.WishlistItem{
		ID: fmt.Sprintf("%d", time.Now().UnixNano()),
		CountryName: countryName,
		Note: "",
		Status: "Planned",
		CreatedAt: time.Now().UTC(),
	}

	storeMu.Lock()
	defer storeMu.Unlock()
	store[username] = append(store[username], items)
	return items, nil
}

// UpdateWishlistItem updates note and status for a specific item
func UpdateWishlistItem(username, id, note, status string) error {
	if !models.AllowedStatuses[status] {
		return fmt.Errorf("invalid status %q: must be Planned or Visited", status)
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	items, ok := store[username]
	if !ok {
		return errors.New("Wishlist not found")
	}

	for i, item := range items {
		if item.ID == id {
			store[username][i].Note = note
			store[username][i].Status = status
			return nil
		}
	}
	return fmt.Errorf("item %q not found", id)
}


// DeleteWishlistItem removes an item from a user's wishlist
func DeleteWishlistItem(username, id string) error {
	storeMu.Lock()
	defer storeMu.Unlock()

	items, ok := store[username]
	if !ok {
		return errors.New("wishlist not found")
	}
	for i, item := range items {
		if item.ID == id {
			store[username] = append(items[:i], items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("item %q not found", id)
}