package services

import (
	"TravelSphere-Mojammel/models"
	"TravelSphere-Mojammel/utils"
)

// popularLocations are well-known coordinates used for the home page
var popularLocations = []struct {
	Lat, Lon float64
}{
	{48.8566, 2.3522},   // Paris
	{40.7128, -74.0060}, // New York
	{-33.8688, 151.2093}, // Sydney
	{41.9028, 12.4964},  // Rome
}

// GetAttractionsByCoords returns attractions near a country's coordinates
func GetAttractionsByCoords(lat, lon float64) ([]models.Attraction, error) {
	if lat == 0 && lon == 0 {
		return []models.Attraction{}, nil
	}
	return utils.FetchAttractions(lat, lon, 10)
}

// GetPopularAttractions returns a mixed list from well-known cities
func GetPopularAttractions() ([]models.Attraction, error) {
	result := make([]models.Attraction, 0)
	for _, loc := range popularLocations {
		items, err := utils.FetchAttractions(loc.Lat, loc.Lon, 2)
		if err != nil {
			continue // non-fatal, skip this location
		}
		result = append(result, items...)
	}
	return result, nil
}