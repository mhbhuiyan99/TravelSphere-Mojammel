package filters

import (
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
)

// AuthFilter protects routes that require a logged-in session
// Redirects to /login if no session is found
func AuthFilter(ctx *context.Context) {
	username := ctx.Input.Session("username")
	if username == nil {
		ctx.Redirect(http.StatusFound, "/login")
	}
}