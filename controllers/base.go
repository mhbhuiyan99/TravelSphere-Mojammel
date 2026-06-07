package controllers

import beego "github.com/beego/beego/v2/server/web"

// BaseController is embedded by all controllers
type BaseController struct {
	beego.Controller
}

// Prepare() runs before every Get()/Post() and sets *shared* template data
func (c *BaseController) Prepare() {
	// Read session username - set by login, empty if not logged in
	username := c.GetSession("username")

	if username != nil {
		c.Data["Username"] = username.(string)
		c.Data["IsLoggedIn"] = true
	} else {
		c.Data["Username"] = ""
		c.Data["IsLoggedIn"] = false
	}

	// Current path - used by nav to highlight active link
	c.Data["CurrentPath"] = c.Ctx.Request.URL.Path
}