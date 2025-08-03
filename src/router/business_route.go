package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func BusinessRoutes(v1 fiber.Router, b service.BusinessService, u service.UserService) {
	businessController := controller.NewBusinessController(b)

	business := v1.Group("/businesses")

	business.Get("/", m.Auth(u, "getBusinesses"), businessController.GetBusinesses)
	business.Post("/", m.Auth(u, "manageBusinesses"), businessController.CreateBusiness)
	business.Get("/:id", m.Auth(u, "getBusinesses"), businessController.GetBusinessByID)
	business.Put("/:id", m.Auth(u, "manageBusinesses"), businessController.UpdateBusiness)
	business.Delete("/:id", m.Auth(u, "manageBusinesses"), businessController.DeleteBusiness)
}
