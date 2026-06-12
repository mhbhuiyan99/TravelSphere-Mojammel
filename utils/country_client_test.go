package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func mockCountryServer(t *testing.T, data interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(data)
	}))
}

func TestTransformCountries(t *testing.T) {
	tests := []struct {
		name     string
		input    []countryAPIResponse
		wantName string
		wantSlug string
		wantCap  string
		wantLat  float64
		wantLon  float64
	}{
		{
			name: "normal country with coords",
			input: []countryAPIResponse{{
				Names:    struct{ Common string `json:"common"` }{Common: "United States"},
				Capitals: []struct{ Name string `json:"name"` }{{Name: "Washington D.C."}},
				Population: 331000000,
				Region:     "Americas",
				Coordinates: struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				}{Lat: 38.0, Lng: -97.0},
			}},
			wantName: "United States",
			wantSlug: "united-states",
			wantCap:  "Washington D.C.",
			wantLat:  38.0,
			wantLon:  -97.0,
		},
		{
			name: "country without coords defaults to zero",
			input: []countryAPIResponse{{
				Names:    struct{ Common string `json:"common"` }{Common: "Antarctica"},
				Capitals: []struct{ Name string `json:"name"` }{},
			}},
			wantName: "Antarctica",
			wantSlug: "antarctica",
			wantCap:  "",
			wantLat:  0,
			wantLon:  0,
		},
		{
			name: "country name with spaces slugified correctly",
			input: []countryAPIResponse{{
				Names: struct{ Common string `json:"common"` }{Common: "New Zealand"},
				Coordinates: struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				}{Lat: -41.0, Lng: 174.0},
			}},
			wantName: "New Zealand",
			wantSlug: "new-zealand",
			wantLat:  -41.0,
			wantLon:  174.0,
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

func TestToSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase single word", "france", "france"},
		{"uppercase", "FRANCE", "france"},
		{"multi word", "United States", "united-states"},
		{"three words", "Papua New Guinea", "papua-new-guinea"},
		{"already slug", "new-zealand", "new-zealand"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlug(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFirstOrEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"single element", []string{"Dhaka"}, "Dhaka"},
		{"multiple elements returns first", []string{"Paris", "Lyon"}, "Paris"},
		{"empty slice returns empty", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstOrEmpty(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMapValues(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]string
		wantCount int
	}{
		{"single entry", map[string]string{"bn": "Bengali"}, 1},
		{"multiple entries", map[string]string{"en": "English", "fr": "French"}, 2},
		{"empty map", map[string]string{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapValues(tt.input)
			if len(got) != tt.wantCount {
				t.Errorf("got %d values, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestCurrencyNames(t *testing.T) {
	tests := []struct {
		name      string
		input     []currencyType
		wantCount int
	}{
		{"single currency", []currencyType{{Code: "BDT", Name: "Bangladeshi taka"}}, 1},
		{"multiple currencies", []currencyType{{Code: "USD", Name: "US Dollar"}, {Code: "EUR", Name: "Euro"}}, 2},
		{"empty name skipped", []currencyType{{Code: "XXX", Name: ""}}, 0},
		{"empty slice", []currencyType{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := currencyNames(tt.input)
			if len(got) != tt.wantCount {
				t.Errorf("got %d names, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestFetchAllCountriesFromURL(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantCount int
		wantErr   bool
	}{
		{
			name: "valid v5 response returns countries",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := v5Response{}
				resp.Data.Objects = []countryAPIResponse{
					{
						Names:    struct{ Common string `json:"common"` }{Common: "Bangladesh"},
						Capitals: []struct{ Name string `json:"name"` }{{Name: "Dhaka"}},
					},
					{
						Names:    struct{ Common string `json:"common"` }{Common: "France"},
						Capitals: []struct{ Name string `json:"name"` }{{Name: "Paris"}},
					},
				}
				resp.Data.Meta.More = false
				json.NewEncoder(w).Encode(resp)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "empty objects returns empty slice",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := v5Response{}
				resp.Data.Meta.More = false
				json.NewEncoder(w).Encode(resp)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "invalid JSON returns error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("not json"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			got, err := FetchAllCountriesFromURL(server.URL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v, got err: %v", tt.wantErr, err)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("got %d countries, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestFetchAllCountriesFromURL_ServerError(t *testing.T) {
	t.Run("unreachable server returns error", func(t *testing.T) {
		_, err := FetchAllCountriesFromURL("http://127.0.0.1:1")
		if err == nil {
			t.Error("expected error for unreachable server, got nil")
		}
	})
}

func TestRestCountriesBase(t *testing.T) {
	t.Run("returns fallback when config not set", func(t *testing.T) {
		got := restCountriesBase()
		want := "https://api.restcountries.com/countries/v5"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns config value when set", func(t *testing.T) {
		beego.AppConfig.Set("restcountries_base_url", "http://mock-countries.com")
		got := restCountriesBase()
		if got != "http://mock-countries.com" {
			t.Errorf("got %q, want http://mock-countries.com", got)
		}
		beego.AppConfig.Set("restcountries_base_url", "")
	})
}

func TestFetchAllCountries_UsesConfigURL(t *testing.T) {
	t.Run("FetchAllCountries calls through to real URL function", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := v5Response{}
			resp.Data.Objects = []countryAPIResponse{{
				Names:    struct{ Common string `json:"common"` }{Common: "Bangladesh"},
				Capitals: []struct{ Name string `json:"name"` }{{Name: "Dhaka"}},
				Coordinates: struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				}{Lat: 24.0, Lng: 90.0},
			}}
			resp.Data.Meta.More = false
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		beego.AppConfig.Set("restcountries_base_url", server.URL)

		got, err := FetchAllCountries()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("got %d countries, want 1", len(got))
		}
		if got[0].Name != "Bangladesh" {
			t.Errorf("got %q, want Bangladesh", got[0].Name)
		}

		beego.AppConfig.Set("restcountries_base_url", "")
	})
}
