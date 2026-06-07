package controllers

type HomeController struct {
	BaseController
}

// Get renders the home page
// @router / [get]
func (c *HomeController) Get() {
	c.Data["Title"] = "TravelSphere"
	c.TplName = "page/home.tpl"
}