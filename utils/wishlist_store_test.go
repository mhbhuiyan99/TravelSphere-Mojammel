package utils

import (
	"os"
	"testing"

	"TravelSphere-Mojammel/models"
	"time"
)

func TestReadWriteStore(t *testing.T) {
	// Use a temp file so tests don't touch real data
	tmp, _ := os.CreateTemp("", "wishlist_test_*.json")
	tmp.Close()
	defer os.Remove(tmp.Name())

	// Override store path for tests
	originalPath := os.Getenv("wishlist_store_path")
	os.Setenv("wishlist_store_path", tmp.Name())
	defer os.Setenv("wishlist_store_path", originalPath)

	tests := []struct {
		name     string
		input    map[string][]models.WishlistItem
		wantKeys int
	}{
		{
			name:     "write and read back single user",
			wantKeys: 1,
			input: map[string][]models.WishlistItem{
				"mojammel": {
					{ID: "1", CountryName: "Bangladesh", Status: "Planned", CreatedAt: time.Now().UTC()},
				},
			},
		},
		{
			name:     "write and read back multiple users",
			wantKeys: 2,
			input: map[string][]models.WishlistItem{
				"alice": {{ID: "2", CountryName: "France", Status: "Visited", CreatedAt: time.Now().UTC()}},
				"bob":   {{ID: "3", CountryName: "Japan", Status: "Planned", CreatedAt: time.Now().UTC()}},
			},
		},
		{
			name:     "write empty store",
			wantKeys: 0,
			input:    map[string][]models.WishlistItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write
			if err := writeStoreToPath(tmp.Name(), tt.input); err != nil {
				t.Fatalf("WriteStore error: %v", err)
			}
			// Read back
			got, err := readStoreFromPath(tmp.Name())
			if err != nil {
				t.Fatalf("ReadStore error: %v", err)
			}
			if len(got) != tt.wantKeys {
				t.Errorf("got %d keys, want %d", len(got), tt.wantKeys)
			}
		})
	}
}

func TestReadStore_MissingFile(t *testing.T) {
	t.Run("missing file returns empty store not error", func(t *testing.T) {
		got, err := readStoreFromPath("/tmp/does_not_exist_xyz.json")
		if err != nil {
			t.Fatalf("expected no error for missing file, got: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected empty store, got %d keys", len(got))
		}
	})
}