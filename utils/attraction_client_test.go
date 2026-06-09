package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAttractions_Success(t *testing.T) {
	mockBody := otmRadiusResponse{
		Features: []struct {
			Properties otmPlace `json:"properties"`
		}{
			{Properties: otmPlace{Name: "Eiffel Tower", Kinds: "architecture", XID: "x1"}},
			{Properties: otmPlace{Name: "", Kinds: "parks", XID: "x2"}}, // should be skipped
			{Properties: otmPlace{Name: "Louvre", Kinds: "museums", XID: "x3"}},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockBody)
	}))
	defer server.Close()

	tests := []struct {
		name      string
		mockResp  otmRadiusResponse
		wantCount int
		wantFirst string
	}{
		{
			name:      "filters out unnamed places",
			mockResp:  mockBody,
			wantCount: 2,
			wantFirst: "Eiffel Tower",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origBase := otmBase
			_ = origBase 

			results, err := fetchAttractionsFromURL(server.URL, 10)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(results) != tt.wantCount {
				t.Errorf("got %d attractions, want %d", len(results), tt.wantCount)
			}
			if results[0].Name != tt.wantFirst {
				t.Errorf("got first name %q, want %q", results[0].Name, tt.wantFirst)
			}
		})
	}
}