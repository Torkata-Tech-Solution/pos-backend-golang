package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func BusinessUserRoutes(v1 fiber.Router, bu service.BusinessUserService, u service.UserService) {
	businessUserController := controller.NewBusinessUserController(bu)

	businessUser := v1.Group("/business-users")

	businessUser.Get("/", m.Auth(u, "getBusinessUsers"), businessUserController.GetBusinessUsers)
	businessUser.Post("/", m.Auth(u, "manageBusinessUsers"), businessUserController.CreateBusinessUser)
	businessUser.Get("/:id", m.Auth(u, "getBusinessUsers"), businessUserController.GetBusinessUserByID)
	businessUser.Put("/:id", m.Auth(u, "manageBusinessUsers"), businessUserController.UpdateBusinessUser)
	businessUser.Delete("/:id", m.Auth(u, "manageBusinessUsers"), businessUserController.DeleteBusinessUser)

	// Business specific routes
	business := v1.Group("/businesses")
	business.Get("/:businessId/users", m.Auth(u, "getBusinessUsers"), businessUserController.GetBusinessUsersByBusinessID)
}
