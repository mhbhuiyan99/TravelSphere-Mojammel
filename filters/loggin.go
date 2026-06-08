package filters

import (
	"github.com/beego/beego/v2/server/web/context"

	"time"
	"fmt"
)

// LogginFilter logs method, path and duration for every request
func LogginFilter(ctx *context.Context) {
	ctx.Input.SetData("requestStart", time.Now())
}

// LogginAfterFilter logs duration after the response is sent
func LogginAfterFilter(ctx *context.Context) {
	start, ok := ctx.Input.GetData("requestStart").(time.Time)
	if !ok {
		return
	}

	fmt.Printf("[%s] %s - %v\n",
		ctx.Request.Method,
		ctx.Request.URL.Path,
		time.Since(start),
	)
}