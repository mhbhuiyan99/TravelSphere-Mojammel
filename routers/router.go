package routers

import (
	"TravelSphere-Mojammel/controllers"
	"TravelSphere-Mojammel/filters"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	// Auth filter - protects wishlist and dashboard
	beego.InsertFilter("/wishlist", beego.BeforeExec, filters.AuthFilter)
	beego.InsertFilter("/dashboard", beego.BeforeExec, filters.AuthFilter)

	// SSR page routes
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.AuthController{})
	beego.Router("/logout", &controllers.LogoutController{})
}