package handler

import (
	"pixie/handler/controller/example"
	"pixie/handler/middleware/routes_middleware"
	mw "pixie/lib/http/middleware"
	"pixie/lib/http/router"
)

// Обертка для объявления роутов через r.AppendRoute()
func RegisterRoutes(r router.IRouter) {
	// Example
	r.AppendRoute("GET", "/example", example.List, []mw.Middleware{
		routes_middleware.NewExampleMiddleware(),
	})
	r.AppendRoute("GET", "/example/{index}", example.Get, []mw.Middleware{})
	r.AppendRoute("POST", "/example", example.Post, []mw.Middleware{})
}
