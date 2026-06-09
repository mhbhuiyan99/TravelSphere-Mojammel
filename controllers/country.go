package controllers

import "TravelSphere-Mojammel/services"

type CountryController struct {
	BaseController
}

// Get renders the Country Explorer page with a default country list
// @router /countries [get]
func (c *CountryController) Get() {
	countries, err := services.GetAllCountries("", "")
	if err != nil {
		c.Data["Error"] = "Could not load countries. Please try again later."
		c.Data["Countries"] = nil
	} else {
		c.Data["Countries"] = countries
	}

	c.Data["Title"] = "Country Explorer"
	c.TplName = "pages/countries.tpl"
}

// Detail renders the destination detail page
// @router /countries/:slug [get]
func (c *CountryController) Detail() {
	slug := c.Ctx.Input.Param(":slug")

	country, err := services.GetCountryBySlug(slug)
	if err != nil || country == nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["Title"] = "Not Found"
		c.TplName = "pages/404.tpl"
		return
	}

	attractions, err := services.GetAttractionsByCoords(country.Lat, country.Lon)
	if err != nil {
		attractions = nil 
	}

	c.Data["Country"] = country
	c.Data["Attractions"] = attractions
	c.Data["Title"] = country.Name
	c.TplName = "pages/destination.tpl"
}