package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func CustomerRoutes(v1 fiber.Router, c service.CustomerService, u service.UserService) {
	customerController := controller.NewCustomerController(c)

	customer := v1.Group("/customers")

	customer.Get("/", m.Auth(u, "getCustomers"), customerController.GetCustomers)
	customer.Post("/", m.Auth(u, "manageCustomers"), customerController.CreateCustomer)
	customer.Get("/:id", m.Auth(u, "getCustomers"), customerController.GetCustomerByID)
	customer.Put("/:id", m.Auth(u, "manageCustomers"), customerController.UpdateCustomer)
	customer.Delete("/:id", m.Auth(u, "manageCustomers"), customerController.DeleteCustomer)
}
