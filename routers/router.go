package routers

import (
	"TravelSphere-Mojammel/controllers"
	"TravelSphere-Mojammel/controllers/api"
	"TravelSphere-Mojammel/filters"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	// Loggin filter - global, runs on every request
	beego.InsertFilter("/*", beego.BeforeExec, filters.LogginFilter)
	beego.InsertFilter("/*", beego.AfterExec, filters.LogginAfterFilter)

	// Auth filter - protects wishlist and dashboard
	beego.InsertFilter("/wishlist", beego.BeforeExec, filters.AuthFilter)
	beego.InsertFilter("/dashboard", beego.BeforeExec, filters.AuthFilter)

	// SSR page routes
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.AuthController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/countries", &controllers.CountryController{})

    // JSON API routes
    beego.Router("/api/countries", &api.CountryAPIController{}, "get:GetAll")
}