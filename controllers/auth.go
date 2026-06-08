package controllers

type AuthController struct {
	BaseController
}

// Get renders the login page
// @router /login [get]
func (c *AuthController) Get() {
	// Already logged in, redirect to home page
	if c.Data["IsLoggedIn"].(bool) {
		c.Redirect("/", 302)
		return
	}
	c.Data["Title"] = "Login"
	c.TplName = "pages/login.tpl"
}

// Post handles the login form submission
// @router /login [post]
func (c *AuthController) Post() {
	name := c.GetString("username")
	
	if name == "" {
		c.Data["Error"] = "Please enter your name"
		c.Data["Title"] = "Login"
		c.TplName = "pages/login.tpl"
		return
	}

	c.SetSession("username", name)
	c.Redirect("/", 302)
}