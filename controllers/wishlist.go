package controllers

import "TravelSphere-Mojammel/services"

type WishlistController struct {
	BaseController
}

// Get renders the wishlist page — auth filter already guards this route
// @router /wishlist [get]
func (c *WishlistController) Get() {
	username := c.Data["Username"].(string)
	items := services.GetWishlist(username)
	c.Data["Items"] = items
	c.Data["Title"] = "Travel Wishlist"
	c.TplName = "pages/wishlist.tpl"
}