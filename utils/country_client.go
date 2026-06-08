package utils

import (
	"TravelSphere-Mojammel/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const restCountriesBase = "https://restcountries.com/v3.1"

type currencyInfo struct {
	Name string `json:"name"`
}

// countryAPIResponse maps the REST Countries JSON shape
type countryAPIResponse struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Capital    []string `json:"capital"`
	Population int64    `json:"population"`
	Region     string   `json:"region"`
	Subregion  string   `json:"subregion"`
	Flags      struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"flags"`
	Languages  map[string]string `json:"languages"`
	Currencies map[string]currencyInfo `json:"currencies"`
}

// FetchAllCountries fetches every country from REST Countries API
func FetchAllCountries() ([]models.Country, error) {
	url := fmt.Sprintf("%s/all?fields=name,capital,population,region,subregion,flags,languages,currencies", restCountriesBase)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("country API request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading country API response: %w", err)
	}

	var raw []countryAPIResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing country API response: %w", err)
	}

	return transformCountries(raw), nil
}

// FetchCountryByName fetches a single country by common name
func FetchCountryByName(name string) (*models.Country, error) {
	url := fmt.Sprintf("%s/name/%s?fullText=true&fields=name,capital,population,region,subregion,flags,languages,currencies", restCountriesBase, name)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("country API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var raw []countryAPIResponse
	if err := json.Unmarshal(body, &raw); err != nil || len(raw) == 0 {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	countries := transformCountries(raw)
	return &countries[0], nil
}

// transformCountries converts raw API response into our domain models
func transformCountries(raw []countryAPIResponse) []models.Country {
	countries := make([]models.Country, 0, len(raw))
	for _, r := range raw {
		countries = append(countries, models.Country{
			Name:       r.Name.Common,
			Slug:       toSlug(r.Name.Common),
			Capital:    firstOrEmpty(r.Capital),
			Population: r.Population,
			Region:     r.Region,
			Subregion:  r.Subregion,
			Flag:       r.Flags.SVG,
			Languages:  mapValues(r.Languages),
			Currencies: currencyNames(r.Currencies),
		})
	}
	return countries
}

func toSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func firstOrEmpty(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

func mapValues(m map[string]string) []string {
	val := make([]string, 0, len(m))

	for _, v := range m {
		val = append(val, v)
	}
	return val
}

func currencyNames(m map[string]currencyInfo) []string {
    names := make([]string, 0, len(m))
    for _, n := range m {
        names = append(names, n.Name)
    }
    return names
}
