package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"TravelSphere-Mojammel/models"

	beego "github.com/beego/beego/v2/server/web"
)

var fileMu sync.RWMutex

func storePath() string {
	path, _ := beego.AppConfig.String("wishlist_store_path")
	if path == "" {
		return "data/wishlist.json"
	}
	return path
}

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

func writeStoreToPath(path string, store map[string][]models.WishlistItem) error {
	fileMu.Lock()
	defer fileMu.Unlock()

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