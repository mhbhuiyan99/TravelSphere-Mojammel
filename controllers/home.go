package controllers

import "TravelSphere-Mojammel/services"

// HomeController handles GET /
type HomeController struct {
	BaseController
}

// Get renders the home page with featured countries and popular attractions
// @router / [get]
func (c *HomeController) Get() {
	// Featured countries — hardcoded slugs, real data from API
	featured, err := services.GetCountriesBySlugs([]string{
		"afghanistan", "russia", "syria", "united-states", "france", "japan", "australia", "bangladesh",
	})
	if err != nil || featured == nil {
		featured = nil
	}

	attractions, err := services.GetPopularAttractions()
	if err != nil {
		attractions = nil
	}

	c.Data["Featured"] = featured
	c.Data["Attractions"] = attractions
	c.Data["Title"] = "TravelSphere"
	c.TplName = "pages/home.tpl"
}