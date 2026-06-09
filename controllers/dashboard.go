package controllers

import "TravelSphere-Mojammel/services"

type DashboardController struct {
	BaseController
}

// Get renders the dashboard page
// @router /dashboard [get]
func (c *DashboardController) Get() {
	username := c.Data["Username"].(string)
	items := services.GetWishlist(username)

	planned, visited := 0, 0
	for _, item := range items {
		if item.Status == "Visited" {
			visited++
		} else {
			planned++
		}
	}

	c.Data["Items"] = items
	c.Data["Total"] = len(items)
	c.Data["Planned"] = planned
	c.Data["Visited"] = visited
	c.Data["Title"] = "Dashboard"
	c.TplName = "pages/dashboard.tpl"
}