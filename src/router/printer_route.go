package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func PrinterRoutes(v1 fiber.Router, p service.PrinterService, u service.UserService) {
	printerController := controller.NewPrinterController(p)

	printer := v1.Group("/printers")

	printer.Get("/", m.Auth(u, "getPrinters"), printerController.GetPrinters)
	printer.Post("/", m.Auth(u, "managePrinters"), printerController.CreatePrinter)
	printer.Get("/:id", m.Auth(u, "getPrinters"), printerController.GetPrinterByID)
	printer.Put("/:id", m.Auth(u, "managePrinters"), printerController.UpdatePrinter)
	printer.Delete("/:id", m.Auth(u, "managePrinters"), printerController.DeletePrinter)

	// Outlet specific routes
	outlet := v1.Group("/outlets")
	outlet.Get("/:outletId/printers", m.Auth(u, "getPrinters"), printerController.GetPrintersByOutletID)
}
