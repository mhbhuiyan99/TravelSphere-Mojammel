package services

import (
	"TravelSphere-Mojammel/models"
	"TravelSphere-Mojammel/utils"
)

//GetAttractionByCoords returns attractions near a country's coordinates
func GetAttractionsByCoords(lat, lon float64) ([]models.Attraction, error) {
	if lat == 0 && lon == 0 {
		return []models.Attraction{}, nil
	}
	return utils.FetchAttractions(lat, lon, 10)
}