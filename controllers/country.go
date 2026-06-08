package controllers

import "TravelSphere-Mojammel/services"

type ContryController struct {
	BaseController
}

// Get renders the Country Explorer page with a default country list
// @router /countries [get]
func (c *ContryController) Get() {
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