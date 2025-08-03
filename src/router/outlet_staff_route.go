package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func OutletStaffRoutes(v1 fiber.Router, os service.OutletStaffService, u service.UserService) {
	outletStaffController := controller.NewOutletStaffController(os)

	outletStaff := v1.Group("/outlet-staff")

	outletStaff.Get("/", m.Auth(u, "getOutletStaff"), outletStaffController.GetOutletStaff)
	outletStaff.Post("/", m.Auth(u, "manageOutletStaff"), outletStaffController.CreateOutletStaff)
	outletStaff.Get("/:id", m.Auth(u, "getOutletStaff"), outletStaffController.GetOutletStaffByID)
	outletStaff.Put("/:id", m.Auth(u, "manageOutletStaff"), outletStaffController.UpdateOutletStaff)
	outletStaff.Delete("/:id", m.Auth(u, "manageOutletStaff"), outletStaffController.DeleteOutletStaff)

	// Outlet specific routes
	outlet := v1.Group("/outlets")
	outlet.Get("/:outletId/staff", m.Auth(u, "getOutletStaff"), outletStaffController.GetOutletStaffByOutletID)
}
