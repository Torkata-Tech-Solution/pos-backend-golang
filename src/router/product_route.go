package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(v1 fiber.Router, p service.ProductService, pc service.ProductCategoryService, u service.UserService) {
	productController := controller.NewProductController(p, pc)

	// Product Category routes
	productCategory := v1.Group("/product-categories")
	productCategory.Get("/", m.Auth(u, "getProductCategories"), productController.GetProductCategories)
	productCategory.Post("/", m.Auth(u, "manageProductCategories"), productController.CreateProductCategory)
	productCategory.Get("/:id", m.Auth(u, "getProductCategories"), productController.GetProductCategoryByID)
	productCategory.Put("/:id", m.Auth(u, "manageProductCategories"), productController.UpdateProductCategory)
	productCategory.Delete("/:id", m.Auth(u, "manageProductCategories"), productController.DeleteProductCategory)

	// Product routes
	product := v1.Group("/products")
	product.Get("/", m.Auth(u, "getProducts"), productController.GetProducts)
	product.Post("/", m.Auth(u, "manageProducts"), productController.CreateProduct)
	product.Get("/:id", m.Auth(u, "getProducts"), productController.GetProductByID)
	product.Put("/:id", m.Auth(u, "manageProducts"), productController.UpdateProduct)
	product.Delete("/:id", m.Auth(u, "manageProducts"), productController.DeleteProduct)
}
