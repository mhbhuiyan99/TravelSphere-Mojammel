package main

import (
	_ "TravelSphere-Mojammel/routers"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {

	// Register template functions: remove [] from slice
	beego.AddFuncMap("join", func(s []string, sep string) string {
		return strings.Join(s, sep)
	})

	beego.Run()
}
