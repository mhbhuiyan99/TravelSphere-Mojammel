package api

import (
	"TravelSphere-Mojammel/services"

	beego "github.com/beego/beego/v2/server/web"
)

type CountryAPIController struct {
	beego.Controller
}

// GetAll returns a filtered JSON list of countries
// @router /api/countriesi [get]
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
		"data": countries,
	}
	c.ServeJSON()
}