package utils

import "testing"

func TestTransformCountries(t *testing.T) {
	tests := []struct {
		name string
		input []countryAPIResponse
		wantName string
		wantSlug string
		wantCap string
		wantLat float64
		wantLon float64
	} {
		{
			name: "normal country with coords",
			input: []countryAPIResponse {
				{
					Name: struct{ Common string `json:"common"` }{Common: "United States"},
					Capital: []string{"Washington D.C."},
					Population: 331000000,
					Region:     "Americas",
					Latlng:     []float64{38.0, -97.0},
				},
			},
			wantName: "United States",
			wantSlug: "united-states",
			wantCap: "Washington D.C.",
			wantLat: 38.0,
			wantLon: -97.0,
		},
		{
			name: "country without coords",
			input: []countryAPIResponse{
				{
					Name:    struct{ Common string `json:"common"` }{Common: "Antarctica"},
					Capital: []string{},
					Latlng:  []float64{},
				},
			},
			wantName: "Antarctica",
			wantSlug: "antarctica",
			wantCap:  "",
			wantLat:  0,
			wantLon:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformCountries(tt.input)
			if len(got) != 1 {
				t.Fatalf("expected 1 country, got %d", len(got))
			}
			c := got[0]
			if c.Name != tt.wantName {
				t.Errorf("Name: got %q want %q", c.Name, tt.wantName)
			}
			if c.Slug != tt.wantSlug {
				t.Errorf("Slug: got %q want %q", c.Slug, tt.wantSlug)
			}
			if c.Capital != tt.wantCap {
				t.Errorf("Capital: got %q want %q", c.Capital, tt.wantCap)
			}
			if c.Lat != tt.wantLat {
				t.Errorf("Lat: got %f want %f", c.Lat, tt.wantLat)
			}
			if c.Lon != tt.wantLon {
				t.Errorf("Lon: got %f want %f", c.Lon, tt.wantLon)
			}
		})
	}
}

/*
func TestFetchAllCountries_MockServer(t *testing.T) {
	mockData := []countryAPIResponse{
		{
			Name:       struct{ Common string `json:"common"` }{Common: "Bangladesh"},
			Capital:    []string{"Dhaka"},
			Population: 170000000,
			Region:     "Asia",
			Latlng:     []float64{24.0, 90.0},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	countries, err := fetchAllCountriesFromURL(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 {
		t.Fatalf("expected 1 country, got %d", len(countries))
	}
	if countries[0].Name != "Bangladesh" {
		t.Errorf("got %q, want Bangladesh", countries[0].Name)
	}
}
*/