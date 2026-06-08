package main

import (
	_ "TravelSphere-Mojammel/routers"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

func main() {

	// Register template functions: remove [] from slice
	beego.AddFuncMap("join", func(s []string, sep string) string {
		return strings.Join(s, sep)
	})

	beego.Run()
}

