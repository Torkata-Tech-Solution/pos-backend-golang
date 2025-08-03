package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func SettingRoutes(v1 fiber.Router, s service.SettingService, u service.UserService) {
	settingController := controller.NewSettingController(s)

	setting := v1.Group("/settings")

	setting.Get("/", m.Auth(u, "getSettings"), settingController.GetSettings)
	setting.Post("/", m.Auth(u, "manageSettings"), settingController.CreateSetting)
	setting.Get("/:id", m.Auth(u, "getSettings"), settingController.GetSettingByID)
	setting.Put("/:id", m.Auth(u, "manageSettings"), settingController.UpdateSetting)
	setting.Delete("/:id", m.Auth(u, "manageSettings"), settingController.DeleteSetting)

	// Get by key route
	setting.Get("/key/:key", m.Auth(u, "getSettings"), settingController.GetSettingByKey)
}
