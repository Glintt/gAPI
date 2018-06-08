package controllers

import (
	"gAPIManagement/api/cache"
	"github.com/qiangxue/fasthttp-routing"
)



func InvalidateCache(c *routing.Context) error {
	cache.InvalidateCache()
	c.Response.SetBody([]byte(`{"error":false, "msg": "Invalidation finished."}`))
	c.Response.Header.SetContentType("application/json")
	return nil
}

func ReloadServices(c *routing.Context) error {
	// InitServices()

	cache.InvalidateCache()
	
	c.Response.SetBody([]byte(`{"error":false, "msg": "Reloaded successfully."}`))
	c.Response.Header.SetContentType("application/json")
	return nil
}