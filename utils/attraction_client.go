package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"TravelSphere-Mojammel/models"
)

func otmBase() string {
    return configOrDefault("opentripmap_base_url", "https://api.opentripmap.com/0.1/en/places")
}

type otmPlace struct {
	Name  string `json:"name"`
	Kinds string `json:"kinds"`
	XID   string `json:"xid"`
}

type otmRadiusResponse struct {
	Features []struct {
		Properties otmPlace `json:"properties"`
	} `json:"features"`
}

func FetchAttractions(lat, lon float64, limit int) ([]models.Attraction, error) {
    apiKey := configOrDefault("opentripmap_api_key", "")
    url := fmt.Sprintf(
        "%s/radius?radius=10000&lon=%f&lat=%f&limit=%d&apikey=%s",
        otmBase(), lon, lat, limit, apiKey,
    )
    return fetchAttractionsFromURL(url, limit)
}


func fetchAttractionsFromURL(url string, limit int) ([]models.Attraction, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("attraction API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading attraction response: %w", err)
	}

	var raw otmRadiusResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing attraction response: %w", err)
	}

	attractions := make([]models.Attraction, 0)
	for _, f := range raw.Features {
		if f.Properties.Name == "" {
			continue
		}
		attractions = append(attractions, models.Attraction{
			Name:  f.Properties.Name,
			Kinds: f.Properties.Kinds,
			XID:   f.Properties.XID,
		})
	}
	return attractions, nil
}
