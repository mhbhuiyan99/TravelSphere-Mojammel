package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"TravelSphere-Mojammel/models"
	"TravelSphere-Mojammel/utils"
)

// v5MockResponse builds a properly wrapped v5 response for mock servers
func v5MockResponse(countries []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"objects": countries,
			"meta": map[string]interface{}{
				"more": false,
			},
		},
	}
}

func mockCountryHandler(countries []map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(v5MockResponse(countries))
	}
}

func TestGetCountriesBySlugs(t *testing.T) {
	tests := []struct {
		name      string
		rawData   []map[string]interface{}
		slugs     []string
		wantCount int
	}{
		{
			name: "two matching slugs",
			rawData: []map[string]interface{}{
				{"names": map[string]string{"common": "France"}, "region": "Europe", "capitals": []map[string]string{{"name": "Paris"}}, "coordinates": map[string]float64{"lat": 46.0, "lng": 2.0}},
				{"names": map[string]string{"common": "Japan"}, "region": "Asia", "capitals": []map[string]string{{"name": "Tokyo"}}, "coordinates": map[string]float64{"lat": 36.0, "lng": 138.0}},
				{"names": map[string]string{"common": "Brazil"}, "region": "Americas", "capitals": []map[string]string{{"name": "Brasilia"}}, "coordinates": map[string]float64{"lat": -10.0, "lng": -51.0}},
			},
			slugs:     []string{"france", "japan"},
			wantCount: 2,
		},
		{
			name: "no matching slugs returns empty",
			rawData: []map[string]interface{}{
				{"names": map[string]string{"common": "France"}, "coordinates": map[string]float64{"lat": 46.0, "lng": 2.0}},
			},
			slugs:     []string{"zzz-unknown"},
			wantCount: 0,
		},
		{
			name: "all slugs match",
			rawData: []map[string]interface{}{
				{"names": map[string]string{"common": "France"}, "coordinates": map[string]float64{"lat": 46.0, "lng": 2.0}},
				{"names": map[string]string{"common": "Japan"}, "coordinates": map[string]float64{"lat": 36.0, "lng": 138.0}},
				{"names": map[string]string{"common": "Brazil"}, "coordinates": map[string]float64{"lat": -10.0, "lng": -51.0}},
			},
			slugs:     []string{"france", "japan", "brazil"},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(mockCountryHandler(tt.rawData))
			defer server.Close()

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

func TestGetCountryBySlug(t *testing.T) {
	tests := []struct {
		name      string
		rawData   []map[string]interface{}
		slug      string
		wantFound bool
		wantName  string
	}{
		{
			name: "existing slug returns country",
			rawData: []map[string]interface{}{
				{"names": map[string]string{"common": "Bangladesh"}, "capitals": []map[string]string{{"name": "Dhaka"}}, "coordinates": map[string]float64{"lat": 24.0, "lng": 90.0}},
			},
			slug:      "bangladesh",
			wantFound: true,
			wantName:  "Bangladesh",
		},
		{
			name: "unknown slug returns nil",
			rawData: []map[string]interface{}{
				{"names": map[string]string{"common": "Bangladesh"}, "coordinates": map[string]float64{"lat": 24.0, "lng": 90.0}},
			},
			slug:      "zzz-unknown",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(mockCountryHandler(tt.rawData))
			defer server.Close()

			all, err := utils.FetchAllCountriesFromURL(server.URL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var found *models.Country
			for _, c := range all {
				if c.Slug == tt.slug {
					found = &c
					break
				}
			}

			if tt.wantFound && found == nil {
				t.Errorf("expected to find slug %q, got nil", tt.slug)
			}
			if !tt.wantFound && found != nil {
				t.Errorf("expected nil for slug %q, got %+v", tt.slug, found)
			}
			if tt.wantFound && found != nil && found.Name != tt.wantName {
				t.Errorf("got name %q, want %q", found.Name, tt.wantName)
			}
		})
	}
}

func TestFilterCountries(t *testing.T) {
	countries := []models.Country{
		{Name: "Bangladesh", Capital: "Dhaka", Region: "Asia"},
		{Name: "France", Capital: "Paris", Region: "Europe"},
		{Name: "Japan", Capital: "Tokyo", Region: "Asia"},
	}

	tests := []struct {
		name      string
		search    string
		region    string
		wantCount int
	}{
		{"no filter returns all", "", "", 3},
		{"search by name", "bang", "", 1},
		{"search by capital", "tokyo", "", 1},
		{"filter by region", "", "Asia", 2},
		{"search and region combined", "france", "Europe", 1},
		{"no match returns empty", "zzz", "", 0},
		{"case insensitive search", "FRANCE", "", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterCountries(countries, tt.search, tt.region)
			if len(got) != tt.wantCount {
				t.Errorf("got %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}