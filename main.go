package main

import (
	_ "TravelSphere-Mojammel/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

