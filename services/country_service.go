package services

import (
	"TravelSphere-Mojammel/models"
	"TravelSphere-Mojammel/utils"
	"sort"
	"strings"
)

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

// GetCountriesBySlugs returns a list of countries matching the given slugs
func GetCountriesBySlugs(slugs []string) ([]models.Country, error) {
	all, err := utils.FetchAllCountries()
	if err != nil {
		return nil, err
	}

	slugSet := make(map[string]bool, len(slugs))
	for _, s := range slugs {
		slugSet[s] = true
	}

	result := make([]models.Country, 0, len(slugs))
	for _, c := range all {
		if slugSet[c.Slug] {
			result = append(result, c)
		}
	}
	return result, nil
}

// filterCountries is the testable core for search/region filtering
func filterCountries(countries []models.Country, search, region string) []models.Country {
	if search == "" && region == "" {
		return countries
	}
	filtered := make([]models.Country, 0)
	for _, c := range countries {
		matchSearch := search == "" ||
			strings.Contains(strings.ToLower(c.Name), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(c.Capital), strings.ToLower(search))
		matchRegion := region == "" || strings.EqualFold(c.Region, region)
		if matchSearch && matchRegion {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// GetAllCountries fetches all countries, optionally filtered
func GetAllCountries(search, region string) ([]models.Country, error) {
	countries, err := utils.FetchAllCountries()
	if err != nil {
		return nil, err
	}
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})
	return filterCountries(countries, search, region), nil
}

func getCountriesBySlugFromURL(url string, slugs []string) ([]models.Country, error) {
	all, err := utils.FetchAllCountriesFromURL(url)
	if err != nil {
		return nil, err
	}
	slugSet := make(map[string]bool)
	for _, s := range slugs {
		slugSet[s] = true
	}
	result := make([]models.Country, 0)
	for _, c := range all {
		if slugSet[c.Slug] {
			result = append(result, c)
		}
	}
	return result, nil
}