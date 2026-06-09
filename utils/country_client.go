package utils

import (
	"TravelSphere-Mojammel/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func restCountriesBase() string {
    return configOrDefault("restcountries_base_url", "https://restcountries.com/v3.1")
}

type currencyType struct {
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
	Currencies map[string]currencyType `json:"currencies"`
	Latlng []float64 `json:"latlng"`
}

func fetchAllCountriesFromURL(url string) ([]models.Country, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var raw []countryAPIResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return transformCountries(raw), nil
}

// FetchAllCountries is the public function used by services
func FetchAllCountries() ([]models.Country, error) {
    url := fmt.Sprintf("%s/all?fields=name,capital,population,region,subregion,flags,languages,currencies,latlng", restCountriesBase())
    return fetchAllCountriesFromURL(url)
}

// FetchCountryByName fetches a single country by common name
func FetchCountryByName(name string) (*models.Country, error) {
	url := fmt.Sprintf("%s/name/%s?fullText=true&fields=name,capital,population,region,subregion,flags,languages,currencies,latlng", restCountriesBase(), name)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("country API request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	body, err := io.ReadAll(res.Body)
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
		c := models.Country{
			Name:       r.Name.Common,
			Slug:       toSlug(r.Name.Common),
			Capital:    firstOrEmpty(r.Capital),
			Population: r.Population,
			Region:     r.Region,
			Subregion:  r.Subregion,
			Flag:       r.Flags.SVG,
			Languages:  mapValues(r.Languages),
			Currencies: currencyNames(r.Currencies),
		}
		if len(r.Latlng) == 2 {
			c.Lat = r.Latlng[0]
			c.Lon = r.Latlng[1]
		}

		countries = append(countries, c)
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

func currencyNames(m map[string]currencyType) []string {
    names := make([]string, 0, len(m))
    for _, n := range m {
        names = append(names, n.Name)
    }
    return names
}
