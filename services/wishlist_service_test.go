package services

import (
	"os"
	"testing"
)

// setupTestStore creates a temp file and points the service at it
func setupTestStore(t *testing.T) {
	t.Helper()
	tmp, err := os.CreateTemp("", "wishlist_test_*.json")
	if err != nil {
		t.Fatalf("could not create temp store: %v", err)
	}
	tmp.WriteString("{}")
	tmp.Close()
	testStorePath = tmp.Name()
	t.Cleanup(func() {
		os.Remove(tmp.Name())
		testStorePath = ""
	})
}

func TestAddWishlist(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		countryName string
		wantErr     bool
	}{
		{"valid entry", "mojammel", "Bangladesh", false},
		{"empty username", "", "Bangladesh", true},
		{"empty country", "mojammel", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestStore(t)
			item, err := AddToWishlist(tt.username, tt.countryName)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr %v, got err: %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if item.CountryName != tt.countryName {
					t.Errorf("got country %q, want %q", item.CountryName, tt.countryName)
				}
				if item.Status != "Planned" {
					t.Errorf("default status should be Planned, got %q", item.Status)
				}
				if item.ID == "" {
					t.Error("ID should not be empty")
				}
			}
		})
	}
}

func TestUpdateWishlistItem(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{"valid status Visited", "Visited", false},
		{"valid status Planned", "Planned", false},
		{"invalid status", "Maybe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestStore(t)
			item, _ := AddToWishlist("mojammel", "France")
			err := UpdateWishlistItem("mojammel", item.ID, "test note", tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr %v, got err: %v", tt.wantErr, err)
			}
		})
	}
}

func TestDeleteWishlistItem(t *testing.T) {
	tests := []struct {
		name    string
		useReal bool
		wantErr bool
	}{
		{"delete existing item", true, false},
		{"delete non-existent id", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestStore(t)
			item, _ := AddToWishlist("mojammel", "Japan")

			id := "bad-id"
			if tt.useReal {
				id = item.ID
			}

			err := DeleteWishlistItem("mojammel", id)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr %v, got err: %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetWishlist(t *testing.T) {
	t.Run("returns empty for new user", func(t *testing.T) {
		setupTestStore(t)
		items := GetWishlist("nobody")
		if len(items) != 0 {
			t.Errorf("expected empty, got %d", len(items))
		}
	})

	t.Run("returns added items", func(t *testing.T) {
		setupTestStore(t)
		AddToWishlist("mojammel", "Italy")
		AddToWishlist("mojammel", "Spain")
		items := GetWishlist("mojammel")
		if len(items) != 2 {
			t.Errorf("expected 2 items, got %d", len(items))
		}
	})
}

// TestWishlistIsolation verifies two users have separate wishlists
func TestWishlistIsolation(t *testing.T) {
	t.Run("separate users have separate wishlists", func(t *testing.T) {
		setupTestStore(t)
		AddToWishlist("alice", "France")
		AddToWishlist("bob", "Japan")
		AddToWishlist("bob", "Korea")

		alice := GetWishlist("alice")
		bob := GetWishlist("bob")

		if len(alice) != 1 {
			t.Errorf("alice: expected 1 item, got %d", len(alice))
		}
		if len(bob) != 2 {
			t.Errorf("bob: expected 2 items, got %d", len(bob))
		}
	})
}