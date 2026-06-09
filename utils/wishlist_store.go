package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"path/filepath"

	"TravelSphere-Mojammel/models"

	beego "github.com/beego/beego/v2/server/web"
)

var fileMu sync.RWMutex

func storePath() string {
	val, err := beego.AppConfig.String("wishlist_store_path")
	if err != nil || val == "" {
		return "data/wishlist.json"
	}
	return val
}

// readStoreFromPath is the testable core
func readStoreFromPath(path string) (map[string][]models.WishlistItem, error) {
	fileMu.RLock()
	defer fileMu.RUnlock()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]models.WishlistItem{}, nil
		}
		return nil, fmt.Errorf("reading wishlist store: %w", err)
	}

	var store map[string][]models.WishlistItem
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parsing wishlist store: %w", err)
	}
	return store, nil
}

// writeStoreToPath is the testable core
func writeStoreToPath(path string, store map[string][]models.WishlistItem) error {
	fileMu.Lock()
	defer fileMu.Unlock()

	// ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating store directory: %w", err)
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling wishlist store: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// ReadStore is the public function used by services
func ReadStore() (map[string][]models.WishlistItem, error) {
	return readStoreFromPath(storePath())
}

// WriteStore is the public function used by services
func WriteStore(store map[string][]models.WishlistItem) error {
	return writeStoreToPath(storePath(), store)
}

// ReadStoreFromPath exported for service-level test injection
func ReadStoreFromPath(path string) (map[string][]models.WishlistItem, error) {
	return readStoreFromPath(path)
}

// WriteStoreToPath exported for service-level test injection
func WriteStoreToPath(path string, store map[string][]models.WishlistItem) error {
	return writeStoreToPath(path, store)
}