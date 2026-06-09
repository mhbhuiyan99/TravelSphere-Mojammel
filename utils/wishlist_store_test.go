package utils

import (
	"os"
	"testing"
	"time"

	"TravelSphere-Mojammel/models"
)

func sampleStore() map[string][]models.WishlistItem {
	return map[string][]models.WishlistItem{
		"mojammel": {
			{ID: "1", CountryName: "Bangladesh", Status: "Planned", CreatedAt: time.Now().UTC()},
		},
	}
}

func tempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "store_test_*.json")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestWriteAndReadStore(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string][]models.WishlistItem
		wantKeys int
	}{
		{"single user", sampleStore(), 1},
		{
			"multiple users",
			map[string][]models.WishlistItem{
				"alice": {{ID: "2", CountryName: "France", Status: "Visited", CreatedAt: time.Now().UTC()}},
				"bob":   {{ID: "3", CountryName: "Japan", Status: "Planned", CreatedAt: time.Now().UTC()}},
			},
			2,
		},
		{"empty store", map[string][]models.WishlistItem{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tempFile(t, "{}")
			if err := writeStoreToPath(path, tt.input); err != nil {
				t.Fatalf("write error: %v", err)
			}
			got, err := readStoreFromPath(path)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}
			if len(got) != tt.wantKeys {
				t.Errorf("got %d keys, want %d", len(got), tt.wantKeys)
			}
		})
	}
}

func TestReadStore_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantKeys int
		wantErr  bool
	}{
		{"valid empty object", "{}", 0, false},
		{"missing file returns empty", "", -1, false}, // -1 signals use missing path
		{"invalid JSON returns error", "not json", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.wantKeys == -1 {
				path = "/tmp/does_not_exist_xyz_123.json"
			} else {
				path = tempFile(t, tt.content)
			}

			got, err := readStoreFromPath(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v, got err: %v", tt.wantErr, err)
			}
			if !tt.wantErr && tt.wantKeys >= 0 && len(got) != tt.wantKeys {
				t.Errorf("got %d keys, want %d", len(got), tt.wantKeys)
			}
			// missing file should return empty map not nil
			if tt.wantKeys == -1 && got == nil {
				t.Error("expected empty map for missing file, got nil")
			}
		})
	}
}

func TestWriteStore_CreatesDirectory(t *testing.T) {
	t.Run("creates parent directory if missing", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "store_dir_test_*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		path := dir + "/nested/wishlist.json"
		err = writeStoreToPath(path, map[string][]models.WishlistItem{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("expected file to be created, but it does not exist")
		}
	})
}

func TestRoundTrip_DataIntegrity(t *testing.T) {
	t.Run("data survives write and read unchanged", func(t *testing.T) {
		path := tempFile(t, "{}")
		original := map[string][]models.WishlistItem{
			"mojammel": {
				{ID: "99", CountryName: "Bangladesh", Note: "Visit Dhaka", Status: "Visited"},
			},
		}

		if err := writeStoreToPath(path, original); err != nil {
			t.Fatalf("write error: %v", err)
		}
		got, err := readStoreFromPath(path)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		item := got["mojammel"][0]
		if item.ID != "99" {
			t.Errorf("ID: got %q want 99", item.ID)
		}
		if item.CountryName != "Bangladesh" {
			t.Errorf("CountryName: got %q want Bangladesh", item.CountryName)
		}
		if item.Note != "Visit Dhaka" {
			t.Errorf("Note: got %q want 'Visit Dhaka'", item.Note)
		}
		if item.Status != "Visited" {
			t.Errorf("Status: got %q want Visited", item.Status)
		}
	})
}

func TestPublicWrappers(t *testing.T) {
	t.Run("ReadStoreFromPath exported wrapper works", func(t *testing.T) {
		path := tempFile(t, `{"alice":[]}`)
		got, err := ReadStoreFromPath(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := got["alice"]; !ok {
			t.Error("expected key 'alice' in result")
		}
	})

	t.Run("WriteStoreToPath exported wrapper works", func(t *testing.T) {
		path := tempFile(t, "{}")
		store := map[string][]models.WishlistItem{
			"bob": {{ID: "1", CountryName: "France", Status: "Planned"}},
		}
		if err := WriteStoreToPath(path, store); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, _ := ReadStoreFromPath(path)
		if len(got["bob"]) != 1 {
			t.Errorf("expected 1 item for bob, got %d", len(got["bob"]))
		}
	})

	t.Run("ReadStore uses default path without crashing", func(t *testing.T) {
		// storePath() falls back to "data/wishlist.json" when config missing
		// ReadStore should return empty map not error when file missing
		_, err := ReadStore()
		if err != nil {
			t.Errorf("ReadStore should not error on missing file, got: %v", err)
		}
	})
}