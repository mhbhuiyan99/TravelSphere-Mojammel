package services

import (
	"errors"
	"fmt"
	"time"

	"TravelSphere-Mojammel/models"
	"TravelSphere-Mojammel/utils"
)

// testStorePath allows tests to inject a temp file path.
// Empty string means use the default from config.
var testStorePath = ""

func readStore() (map[string][]models.WishlistItem, error) {
	if testStorePath != "" {
		return utils.ReadStoreFromPath(testStorePath)
	}
	return utils.ReadStore()
}

func writeStore(store map[string][]models.WishlistItem) error {
	if testStorePath != "" {
		return utils.WriteStoreToPath(testStorePath, store)
	}
	return utils.WriteStore(store)
}

// GetWishlist returns all wishlist items for a user
func GetWishlist(username string) []models.WishlistItem {
	store, err := readStore()
	if err != nil {
		return []models.WishlistItem{}
	}
	items, ok := store[username]
	if !ok {
		return []models.WishlistItem{}
	}
	return items
}

// AddToWishlist creates a new wishlist entry and persists it
func AddToWishlist(username, countryName string) (models.WishlistItem, error) {
	if username == "" {
		return models.WishlistItem{}, errors.New("username required")
	}
	if countryName == "" {
		return models.WishlistItem{}, errors.New("country name required")
	}

	store, err := readStore()
	if err != nil {
		return models.WishlistItem{}, fmt.Errorf("loading store: %w", err)
	}

	item := models.WishlistItem{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		CountryName: countryName,
		Note:        "",
		Status:      "Planned",
		CreatedAt:   time.Now().UTC(),
	}

	store[username] = append(store[username], item)

	if err := writeStore(store); err != nil {
		return models.WishlistItem{}, fmt.Errorf("saving store: %w", err)
	}
	return item, nil
}

// UpdateWishlistItem updates note and status for a specific item
func UpdateWishlistItem(username, id, note, status string) error {
	if !models.AllowedStatuses[status] {
		return fmt.Errorf("invalid status %q: must be Planned or Visited", status)
	}

	store, err := readStore()
	if err != nil {
		return fmt.Errorf("loading store: %w", err)
	}

	items, ok := store[username]
	if !ok {
		return errors.New("wishlist not found")
	}
	for i, item := range items {
		if item.ID == id {
			store[username][i].Note = note
			store[username][i].Status = status
			return writeStore(store)
		}
	}
	return fmt.Errorf("item %q not found", id)
}

// DeleteWishlistItem removes an item and persists the change
func DeleteWishlistItem(username, id string) error {
	store, err := readStore()
	if err != nil {
		return fmt.Errorf("loading store: %w", err)
	}

	items, ok := store[username]
	if !ok {
		return errors.New("wishlist not found")
	}
	for i, item := range items {
		if item.ID == id {
			store[username] = append(items[:i], items[i+1:]...)
			return writeStore(store)
		}
	}
	return fmt.Errorf("item %q not found", id)
}