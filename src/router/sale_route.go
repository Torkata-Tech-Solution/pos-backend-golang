package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func SaleRoutes(v1 fiber.Router, s service.SaleService, u service.UserService) {
	saleController := controller.NewSaleController(s)

	sale := v1.Group("/sales")

	sale.Get("/", m.Auth(u, "getSales"), saleController.GetSales)
	sale.Post("/", m.Auth(u, "manageSales"), saleController.CreateSale)
	sale.Get("/:id", m.Auth(u, "getSales"), saleController.GetSaleByID)
	sale.Put("/:id", m.Auth(u, "manageSales"), saleController.UpdateSale)
	sale.Delete("/:id", m.Auth(u, "manageSales"), saleController.DeleteSale)

	// Additional sale routes
	sale.Get("/invoice/:invoiceNumber", m.Auth(u, "getSales"), saleController.GetSaleByInvoiceNumber)
	sale.Patch("/:id/status", m.Auth(u, "manageSales"), saleController.UpdateSaleStatus)
	sale.Get("/report", m.Auth(u, "getSalesReport"), saleController.GetSalesReport)
}
