package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func CouponRoutes(v1 fiber.Router, c service.CouponService, u service.UserService) {
	couponController := controller.NewCouponController(c)

	coupon := v1.Group("/coupons")

	coupon.Get("/", m.Auth(u, "getCoupons"), couponController.GetCoupons)
	coupon.Post("/", m.Auth(u, "manageCoupons"), couponController.CreateCoupon)
	coupon.Get("/:id", m.Auth(u, "getCoupons"), couponController.GetCouponByID)
	coupon.Put("/:id", m.Auth(u, "manageCoupons"), couponController.UpdateCoupon)
	coupon.Delete("/:id", m.Auth(u, "manageCoupons"), couponController.DeleteCoupon)

	// Get by code route
	coupon.Get("/code/:code", m.Auth(u, "getCoupons"), couponController.GetCouponByCode)
}
