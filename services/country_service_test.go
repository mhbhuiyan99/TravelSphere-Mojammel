package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"TravelSphere-Mojammel/models"
)

func TestGetCountriesBySlugs(t *testing.T) {
	mockCountries := []models.Country{
		{Name: "France", Slug: "france", Region: "Europe"},
		{Name: "Japan", Slug: "japan", Region: "Asia"},
		{Name: "Brazil", Slug: "brazil", Region: "Americas"},
	}

	// Mock HTTP server returning raw API shape
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := []map[string]interface{}{
			{"name": map[string]string{"common": "France"}, "region": "Europe", "capital": []string{"Paris"}, "latlng": []float64{46.0, 2.0}},
			{"name": map[string]string{"common": "Japan"}, "region": "Asia", "capital": []string{"Tokyo"}, "latlng": []float64{36.0, 138.0}},
			{"name": map[string]string{"common": "Brazil"}, "region": "Americas", "capital": []string{"Brasilia"}, "latlng": []float64{-10.0, -51.0}},
		}
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	_ = mockCountries // shape reference

	tests := []struct {
		name      string
		slugs     []string
		wantCount int
	}{
		{"two matching slugs", []string{"france", "japan"}, 2},
		{"no matching slugs", []string{"zzz-unknown"}, 0},
		{"all matching", []string{"france", "japan", "brazil"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override fetch URL for test
			got, err := getCountriesBySlugFromURL(server.URL, tt.slugs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("got %d countries, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestGetAllCountries_FilterBySearch(t *testing.T) {
	tests := []struct {
		name      string
		countries []models.Country
		search    string
		region    string
		wantCount int
	}{
		{
			name: "search by name",
			countries: []models.Country{
				{Name: "Bangladesh", Capital: "Dhaka", Region: "Asia"},
				{Name: "France", Capital: "Paris", Region: "Europe"},
			},
			search:    "bang",
			wantCount: 1,
		},
		{
			name: "search by capital",
			countries: []models.Country{
				{Name: "Bangladesh", Capital: "Dhaka", Region: "Asia"},
				{Name: "France", Capital: "Paris", Region: "Europe"},
			},
			search:    "paris",
			wantCount: 1,
		},
		{
			name: "filter by region",
			countries: []models.Country{
				{Name: "Bangladesh", Capital: "Dhaka", Region: "Asia"},
				{Name: "Japan", Capital: "Tokyo", Region: "Asia"},
				{Name: "France", Capital: "Paris", Region: "Europe"},
			},
			region:    "Asia",
			wantCount: 2,
		},
		{
			name: "no match returns empty",
			countries: []models.Country{
				{Name: "France", Capital: "Paris", Region: "Europe"},
			},
			search:    "zzz",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterCountries(tt.countries, tt.search, tt.region)
			if len(got) != tt.wantCount {
				t.Errorf("got %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}