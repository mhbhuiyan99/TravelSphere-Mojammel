package api

import (
	"TravelSphere-Mojammel/services"

	beego "github.com/beego/beego/v2/server/web"
)


type CountryAPIController struct {
	beego.Controller
}

// GetAll returns a filtered JSON list of countries
// @router /api/countries [get]
func (c *CountryAPIController) GetAll() {
	search := c.GetString("search")
	region := c.GetString("region")

	countries, err := services.GetAllCountries(search, region)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Failed to fetch countries",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    countries,
	}
	c.ServeJSON()
}

// GetBySlug returns JSON detail for a single country
// @router /api/countries/:slug [get]
func (c *CountryAPIController) GetBySlug() {
	slug := c.Ctx.Input.Param(":slug")

	country, err := services.GetCountryBySlug(slug)
	if err != nil || country == nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Country not found",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    country,
	}
	c.ServeJSON()
}