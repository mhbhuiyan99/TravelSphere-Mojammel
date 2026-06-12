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
	return configOrDefault("restcountries_base_url", "https://api.restcountries.com/countries/v5")
}

func restCountriesAPIKey() string {
	return configOrDefault("restcountries_api_key", "")
}

type currencyType struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// v5 wraps all responses under data.objects + data.meta
type v5Response struct {
	Data struct {
		Objects []countryAPIResponse `json:"objects"`
		Meta    struct {
			More bool `json:"more"`
		} `json:"meta"`
	} `json:"data"`
}

// countryAPIResponse maps the REST Countries v5 JSON shape
type countryAPIResponse struct {
	Names struct {
		Common string `json:"common"`
	} `json:"names"`
	Capitals []struct {
		Name string `json:"name"`
	} `json:"capitals"`
	Population int64  `json:"population"`
	Region     string `json:"region"`
	Subregion  string `json:"subregion"`
	Flag       struct {
		URLSvg string `json:"url_svg"`
		URLPng string `json:"url_png"`
	} `json:"flag"`
	Languages []struct {
		Name string `json:"name"`
	} `json:"languages"`
	Currencies  []currencyType `json:"currencies"` 
	Coordinates struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinates"`
}

// doAuthedGet makes an authenticated GET to v5 API
func doAuthedGet(url string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	if key := restCountriesAPIKey(); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}
	return body, resp.StatusCode, nil
}

// fetchPageFromURL fetches one page — testable core used by FetchAllCountriesFromURL
func fetchPageFromURL(url string) ([]countryAPIResponse, bool, error) {
	body, status, err := doAuthedGet(url)
	if err != nil {
		return nil, false, err
	}
	if status != http.StatusOK {
		return nil, false, fmt.Errorf("API returned status %d", status)
	}

	var result v5Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, false, fmt.Errorf("parsing response: %w", err)
	}
	return result.Data.Objects, result.Data.Meta.More, nil
}

// FetchAllCountriesFromURL fetches all pages from a base URL — exported for tests
func FetchAllCountriesFromURL(baseURL string) ([]models.Country, error) {
	all := make([]countryAPIResponse, 0)
	limit := 100
	offset := 0
	fields := "response_fields=names.common,capitals,population,region,subregion,flag.url_svg,flag.url_png,languages,currencies,coordinates.lat,coordinates.lng"

	for {
		url := fmt.Sprintf("%s?%s&limit=%d&offset=%d", baseURL, fields, limit, offset)
		objects, more, err := fetchPageFromURL(url)
		if err != nil {
			return nil, err
		}
		all = append(all, objects...)
		if !more {
			break
		}
		offset += limit
	}
	return transformCountries(all), nil
}

// FetchAllCountries is the public function used by services
func FetchAllCountries() ([]models.Country, error) {
	return FetchAllCountriesFromURL(restCountriesBase())
}

// FetchCountryByName fetches a single country by common name
func FetchCountryByName(name string) (*models.Country, error) {
	encoded := strings.ReplaceAll(name, " ", "+")
	url := fmt.Sprintf("%s/names.common/%s", restCountriesBase(), encoded)

	body, status, err := doAuthedGet(url)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", status)
	}

	var result v5Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	if len(result.Data.Objects) == 0 {
		return nil, nil
	}

	countries := transformCountries(result.Data.Objects)
	return &countries[0], nil
}

// transformCountries converts raw v5 objects into domain models
func transformCountries(raw []countryAPIResponse) []models.Country {
	countries := make([]models.Country, 0, len(raw))
	for _, r := range raw {
		c := models.Country{
			Name:       r.Names.Common,
			Slug:       toSlug(r.Names.Common),
			Capital:    firstCapitalName(r.Capitals),
			Population: r.Population,
			Region:     r.Region,
			Subregion:  r.Subregion,
			Flag:       r.Flag.URLSvg,
			Languages:  languageNames(r.Languages),
			Currencies: currencyNames(r.Currencies),
			Lat:        r.Coordinates.Lat,
			Lon:        r.Coordinates.Lng,
		}
		if c.Flag == "" {
			c.Flag = r.Flag.URLPng // fallback to PNG if SVG missing
		}
		countries = append(countries, c)
	}
	return countries
}

func toSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

// firstCapitalName extracts name from the first capital object
func firstCapitalName(capitals []struct {
	Name string `json:"name"`
}) string {
	if len(capitals) > 0 {
		return capitals[0].Name
	}
	return ""
}

// languageNames extracts names from v5 languages array
func languageNames(langs []struct {
	Name string `json:"name"`
}) []string {
	names := make([]string, 0, len(langs))
	for _, l := range langs {
		if l.Name != "" {
			names = append(names, l.Name)
		}
	}
	return names
}

func currencyNames(currencies []currencyType) []string {
	names := make([]string, 0, len(currencies))
	for _, c := range currencies {
		if c.Name != "" {
			names = append(names, c.Name)
		}
	}
	return names
}

// firstOrEmpty kept for any other callers
func firstOrEmpty(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

// mapValues kept for any other callers
func mapValues(m map[string]string) []string {
	val := make([]string, 0, len(m))
	for _, v := range m {
		val = append(val, v)
	}
	return val
}