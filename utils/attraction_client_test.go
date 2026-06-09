package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestFetchAttractionsFromURL(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantCount int
		wantErr   bool
	}{
		{
			name: "filters out unnamed places",
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := otmRadiusResponse{
					Features: []struct {
						Properties otmPlace `json:"properties"`
					}{
						{Properties: otmPlace{Name: "Eiffel Tower", Kinds: "architecture", XID: "x1"}},
						{Properties: otmPlace{Name: "", Kinds: "parks", XID: "x2"}},
						{Properties: otmPlace{Name: "Louvre", Kinds: "museums", XID: "x3"}},
					},
				}
				json.NewEncoder(w).Encode(data)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "all unnamed returns empty slice",
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := otmRadiusResponse{
					Features: []struct {
						Properties otmPlace `json:"properties"`
					}{
						{Properties: otmPlace{Name: "", XID: "x1"}},
						{Properties: otmPlace{Name: "", XID: "x2"}},
					},
				}
				json.NewEncoder(w).Encode(data)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "empty features returns empty slice",
			handler: func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(otmRadiusResponse{})
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "invalid JSON returns error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("bad json"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			got, err := fetchAttractionsFromURL(server.URL, 10)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v, got err: %v", tt.wantErr, err)
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("got %d attractions, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestFetchAttractionsFromURL_NamesPreserved(t *testing.T) {
	t.Run("attraction names and kinds preserved correctly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := otmRadiusResponse{
				Features: []struct {
					Properties otmPlace `json:"properties"`
				}{
					{Properties: otmPlace{Name: "Grand Canyon", Kinds: "natural", XID: "gc1"}},
				},
			}
			json.NewEncoder(w).Encode(data)
		}))
		defer server.Close()

		got, err := fetchAttractionsFromURL(server.URL, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got[0].Name != "Grand Canyon" {
			t.Errorf("got name %q, want Grand Canyon", got[0].Name)
		}
		if got[0].Kinds != "natural" {
			t.Errorf("got kinds %q, want natural", got[0].Kinds)
		}
		if got[0].XID != "gc1" {
			t.Errorf("got XID %q, want gc1", got[0].XID)
		}
	})
}

func TestOtmBase(t *testing.T) {
	t.Run("returns fallback when config not set", func(t *testing.T) {
		got := otmBase()
		want := "https://api.opentripmap.com/0.1/en/places"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns config value when set", func(t *testing.T) {
		beego.AppConfig.Set("opentripmap_base_url", "http://mock-otm.com")
		got := otmBase()
		if got != "http://mock-otm.com" {
			t.Errorf("got %q, want http://mock-otm.com", got)
		}
		// reset
		beego.AppConfig.Set("opentripmap_base_url", "")
	})
}

func TestFetchAttractions_UsesConfigURL(t *testing.T) {
	t.Run("FetchAttractions calls real function path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := otmRadiusResponse{
				Features: []struct {
					Properties otmPlace `json:"properties"`
				}{
					{Properties: otmPlace{Name: "Test Place", Kinds: "historic", XID: "t1"}},
				},
			}
			json.NewEncoder(w).Encode(data)
		}))
		defer server.Close()

		// Point config to mock server
		beego.AppConfig.Set("opentripmap_base_url", server.URL)
		beego.AppConfig.Set("opentripmap_api_key", "testkey")

		got, err := FetchAttractions(0.0, 0.0, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("got %d attractions, want 1", len(got))
		}

		// reset
		beego.AppConfig.Set("opentripmap_base_url", "")
		beego.AppConfig.Set("opentripmap_api_key", "")
	})
}