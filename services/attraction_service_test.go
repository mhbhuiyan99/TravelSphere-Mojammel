package services

import (
	"testing"
)

func TestGetAttractionsByCoords_ZeroCoords(t *testing.T) {
	tests := []struct {
		name      string
		lat, lon  float64
		wantEmpty bool
	}{
		{"zero coords returns empty", 0, 0, true},
		{"valid coords passes through", 48.8566, 2.3522, false}, // Paris — hits real API, skip in CI
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.lat == 0 && tt.lon == 0 {
				result, err := GetAttractionsByCoords(tt.lat, tt.lon)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(result) != 0 {
					t.Errorf("expected empty slice for zero coords, got %d items", len(result))
				}
			}
		})
	}
}