package controllers

type LogoutController struct {
	BaseController
}

// Get clear the session and redirect to login page
// @router /logout [get]
func (c *LogoutController) Get() {
	c.DelSession("username")
	c.Redirect("/login", 302)
}