package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func TableRoutes(v1 fiber.Router, t service.TableService, u service.UserService) {
	tableController := controller.NewTableController(t)

	table := v1.Group("/tables")

	table.Get("/", m.Auth(u, "getTables"), tableController.GetTables)
	table.Post("/", m.Auth(u, "manageTables"), tableController.CreateTable)
	table.Get("/:id", m.Auth(u, "getTables"), tableController.GetTableByID)
	table.Put("/:id", m.Auth(u, "manageTables"), tableController.UpdateTable)
	table.Delete("/:id", m.Auth(u, "manageTables"), tableController.DeleteTable)

	// Outlet specific routes
	outlet := v1.Group("/outlets")
	outlet.Get("/:outletId/tables", m.Auth(u, "getTables"), tableController.GetTablesByOutletID)
}
