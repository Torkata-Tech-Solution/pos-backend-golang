package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func PaymentMethodRoutes(v1 fiber.Router, pm service.PaymentMethodService, u service.UserService) {
	paymentMethodController := controller.NewPaymentMethodController(pm)

	paymentMethod := v1.Group("/payment-methods")

	paymentMethod.Get("/", m.Auth(u, "getPaymentMethods"), paymentMethodController.GetPaymentMethods)
	paymentMethod.Post("/", m.Auth(u, "managePaymentMethods"), paymentMethodController.CreatePaymentMethod)
	paymentMethod.Get("/:id", m.Auth(u, "getPaymentMethods"), paymentMethodController.GetPaymentMethodByID)
	paymentMethod.Put("/:id", m.Auth(u, "managePaymentMethods"), paymentMethodController.UpdatePaymentMethod)
	paymentMethod.Delete("/:id", m.Auth(u, "managePaymentMethods"), paymentMethodController.DeletePaymentMethod)
}
