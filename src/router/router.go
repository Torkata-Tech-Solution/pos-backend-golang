package router

import (
	"app/src/config"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	// Initialize services
	healthCheckService := service.NewHealthCheckService(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(db, validate)
	tokenService := service.NewTokenService(db, validate, userService)
	authService := service.NewAuthService(db, validate, userService, tokenService)
	businessService := service.NewBusinessService(db, validate)
	businessUserService := service.NewBusinessUserService(db, validate)
	outletService := service.NewOutletService(db, validate)
	outletStaffService := service.NewOutletStaffService(db, validate)
	productService := service.NewProductService(db, validate)
	productCategoryService := service.NewProductCategoryService(db, validate)
	customerService := service.NewCustomerService(db, validate)
	tableService := service.NewTableService(db, validate)
	paymentMethodService := service.NewPaymentMethodService(db, validate)
	printerService := service.NewPrinterService(db, validate)
	settingService := service.NewSettingService(db, validate)
	couponService := service.NewCouponService(db, validate)
	saleService := service.NewSaleService(db, validate)

	v1 := app.Group("/v1")

	// Register routes
	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, tokenService, emailService)
	UserRoutes(v1, userService, tokenService)
	BusinessRoutes(v1, businessService, userService)
	BusinessUserRoutes(v1, businessUserService, userService)
	OutletRoutes(v1, outletService, userService)
	OutletStaffRoutes(v1, outletStaffService, userService)
	ProductRoutes(v1, productService, productCategoryService, userService)
	CustomerRoutes(v1, customerService, userService)
	TableRoutes(v1, tableService, userService)
	PaymentMethodRoutes(v1, paymentMethodService, userService)
	PrinterRoutes(v1, printerService, userService)
	SettingRoutes(v1, settingService, userService)
	CouponRoutes(v1, couponService, userService)
	SaleRoutes(v1, saleService, userService)

	if !config.IsProd {
		DocsRoutes(v1)
	}
}
