package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func OutletRoutes(v1 fiber.Router, o service.OutletService, u service.UserService) {
	outletController := controller.NewOutletController(o)

	outlet := v1.Group("/outlets")

	outlet.Get("/", m.Auth(u, "getOutlets"), outletController.GetOutlets)
	outlet.Post("/", m.Auth(u, "manageOutlets"), outletController.CreateOutlet)
	outlet.Get("/:id", m.Auth(u, "getOutlets"), outletController.GetOutletByID)
	outlet.Put("/:id", m.Auth(u, "manageOutlets"), outletController.UpdateOutlet)
	outlet.Delete("/:id", m.Auth(u, "manageOutlets"), outletController.DeleteOutlet)

	// Business specific routes
	business := v1.Group("/businesses")
	business.Get("/:businessId/outlets", m.Auth(u, "getOutlets"), outletController.GetOutletsByBusinessID)
}
