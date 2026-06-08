package services

import (
	"TravelSphere-Mojammel/models"
	"TravelSphere-Mojammel/utils"
	"sort"
	"strings"
)

// GetAllCountries fetches all countries, optionally filtered by search and region
func GetAllCountries(search, region string) ([]models.Country, error) {
	countries, err := utils.FetchAllCountries()
	if err != nil {
		return nil, err
	}	

	// Sort alphabetically by name
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	// Apply filters
	if search == "" && region == "" {
		return countries, nil
	}

	filtered := make([]models.Country, 0)

	for _, c := range countries {
		matchSearch := search == "" || strings.Contains(strings.ToLower(c.Name), strings.ToLower(search)) || 
			strings.Contains(strings.ToLower(c.Capital), strings.ToLower(search))

		matchRegion := region == "" || strings.EqualFold(c.Region, region)

		if matchSearch && matchRegion {
			filtered = append(filtered, c)
		}
	}
	return filtered, nil
}

// GetCountryBySlug finds a single country matching the slug
func GetCountryBySlug(slug string) (*models.Country, error) {
    countries, err := utils.FetchAllCountries()
    if err != nil {
        return nil, err
    }
    for _, c := range countries {
        if c.Slug == slug {
            return &c, nil
        }
    }
    return nil, nil
}