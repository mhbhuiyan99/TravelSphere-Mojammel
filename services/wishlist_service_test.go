package services

import (
	"TravelSphere-Mojammel/models"
	"TravelSphere-Mojammel/utils"
	"testing"
)

func resetStore() {
	utils.WriteStore(map[string][]models.WishlistItem{})
}

func TestAddWishlist(t *testing.T) {
	tests := []struct {
		name string
		username string
		countryName string
		wantErr bool
	}{
		{"valid entry", "mojammel", "Bangladesh", false},
		{"empty username", "", "Bangladesh", true},
		{"empty country", "mojammel", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetStore()
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
			resetStore()
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
		id      string
		wantErr bool
	}{
		{"delete existing item", "", false},   // ID filled in at runtime
		{"delete non-existent id", "bad-id", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetStore()
			item, _ := AddToWishlist("mojammel", "Japan")

			id := tt.id
			if id == "" {
				id = item.ID // use real ID for success case
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
		resetStore()
		items := GetWishlist("nobody")
		if len(items) != 0 {
			t.Errorf("expected empty, got %d", len(items))
		}
	})

	t.Run("returns added items", func(t *testing.T) {
		resetStore()
		AddToWishlist("mojammel", "Italy")
		AddToWishlist("mojammel", "Spain")
		items := GetWishlist("mojammel")
		if len(items) != 2 {
			t.Errorf("expected 2 items, got %d", len(items))
		}
	})
}
