package services

import (
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestGetAttractionsByCoords_ZeroCoords(t *testing.T) {
	tests := []struct {
		name      string
		lat, lon  float64
		wantEmpty bool
	}{
		{"zero coords returns empty without API call", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetAttractionsByCoords(tt.lat, tt.lon)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantEmpty && len(result) != 0 {
				t.Errorf("expected empty, got %d items", len(result))
			}
		})
	}
}

func TestGetPopularAttractions_ReturnsWithoutCrash(t *testing.T) {
	tests := []struct {
		name         string
		mockOTMURL   string
		wantNonNil   bool
		wantNoError  bool
	}{
		{
			name:        "all locations fail gracefully — returns empty non-nil slice",
			mockOTMURL:  "http://127.0.0.1:1", // unreachable
			wantNonNil:  true,
			wantNoError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beego.AppConfig.Set("opentripmap_base_url", tt.mockOTMURL)
			defer beego.AppConfig.Set("opentripmap_base_url", "")

			result, err := GetPopularAttractions()
			if tt.wantNoError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.wantNonNil && result == nil {
				t.Error("expected non-nil slice, got nil")
			}
		})
	}
}