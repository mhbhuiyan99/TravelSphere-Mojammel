package api

import (
	"TravelSphere-Mojammel/services"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type WishlistAPIController struct {
	beego.Controller
}

// username reads the logged-in user from session
func (c *WishlistAPIController) username() string {
	u := c.GetSession("username")
	if u == nil {
		return ""
	}
	return u.(string)
}

// GetAll returns wishlist items for the current user
// @router /api/wishlist [get]
func (c *WishlistAPIController) GetAll() {
	user := c.username()
	if user == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "not authenticated"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    services.GetWishlist(user),
	}
	c.ServeJSON()
}

// Create adds a new wishlist entry
// @router /api/wishlist [post]
func (c *WishlistAPIController) Create() {
	user := c.username()
	if user == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "not authenticated"}
		c.ServeJSON()
		return
	}

	var body struct {
		CountryName string `json:"country_name"`
	}

	// Beego v2 — read raw body and unmarshal manually
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "invalid request body"}
		c.ServeJSON()
		return
	}

	if body.CountryName == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "country_name required"}
		c.ServeJSON()
		return
	}

	item, err := services.AddToWishlist(user, body.CountryName)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "message": err.Error()}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]interface{}{"success": true, "data": item}
	c.ServeJSON()
}

// Update modifies note and status for a wishlist item
// @router /api/wishlist/:id [put]
func (c *WishlistAPIController) Update() {
	user := c.username()
	if user == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "not authenticated"}
		c.ServeJSON()
		return
	}

	id := c.Ctx.Input.Param(":id")

	var body struct {
		Note   string `json:"note"`
		Status string `json:"status"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "invalid request body"}
		c.ServeJSON()
		return
	}

	if err := services.UpdateWishlistItem(user, id, body.Note, body.Status); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "message": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"success": true, "message": "updated"}
	c.ServeJSON()
}

// Delete removes a wishlist item
// @router /api/wishlist/:id [delete]
func (c *WishlistAPIController) Delete() {
	user := c.username()
	if user == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "not authenticated"}
		c.ServeJSON()
		return
	}

	id := c.Ctx.Input.Param(":id")
	if err := services.DeleteWishlistItem(user, id); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"success": false, "message": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"success": true, "message": "deleted"}
	c.ServeJSON()
}

// Summary returns wishlist counts for AJAX dashboard refresh
// @router /api/dashboard/summary [get]
func (c *WishlistAPIController) Summary() {
	user := c.username()
	if user == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "not authenticated"}
		c.ServeJSON()
		return
	}

	items := services.GetWishlist(user)
	planned, visited := 0, 0
	for _, item := range items {
		if item.Status == "Visited" {
			visited++
		} else {
			planned++
		}
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data": map[string]int{
			"total":   len(items),
			"planned": planned,
			"visited": visited,
		},
	}
	c.ServeJSON()
}