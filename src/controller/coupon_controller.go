package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CouponController struct {
	CouponService service.CouponService
}

func NewCouponController(couponService service.CouponService) *CouponController {
	return &CouponController{
		CouponService: couponService,
	}
}

// @Tags         Coupon
// @Summary      Get all coupons
// @Description  Get all coupons with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page      query     int     false   "Page number"  default(1)
// @Param        limit     query     int     false   "Maximum number of coupons"    default(10)
// @Param        search    query     string  false  "Search by code or name"
// @Param        outletId  query     string  false  "Filter by outlet ID"
// @Param        active    query     bool    false  "Filter by active status"
// @Router       /coupons [get]
// @Success      200  {object}  response.SuccessWithPaginatedCoupons
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CouponController) GetCoupons(ctx *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
		Search: ctx.Query("search", ""),
	}

	outletID := ctx.Query("outletId", "")

	coupons, totalResults, err := c.CouponService.GetCoupons(ctx, query, outletID)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedCoupons{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Coupons retrieved successfully",
		Results:      coupons,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Coupon
// @Summary      Get coupon by ID
// @Description  Get a specific coupon by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Coupon ID"
// @Router       /coupons/{id} [get]
// @Success      200  {object}  response.SuccessWithCoupon
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Coupon not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CouponController) GetCouponByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid coupon ID format",
			Errors:  "Coupon ID must be a valid UUID",
		})
	}

	coupon, err := c.CouponService.GetCouponByID(ctx, id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithCoupon{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Coupon retrieved successfully",
		Coupon:  *coupon,
	})
}

// @Tags         Coupon
// @Summary      Get coupon by code
// @Description  Get a specific coupon by its code
// @Security     BearerAuth
// @Produce      json
// @Param        code     path      string  true  "Coupon Code"
// @Param        outletId query     string  false "Outlet ID (optional)"
// @Router       /coupons/code/{code} [get]
// @Success      200  {object}  response.SuccessWithCoupon
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Coupon not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CouponController) GetCouponByCode(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	outletID := ctx.Query("outletId", "")

	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Coupon code is required",
			Errors:  "Coupon code cannot be empty",
		})
	}

	coupon, err := c.CouponService.GetCouponByCode(ctx, code, outletID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithCoupon{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Coupon retrieved successfully",
		Coupon:  *coupon,
	})
}

// @Tags         Coupon
// @Summary      Create new coupon
// @Description  Create a new coupon
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        coupon  body      validation.CreateCoupon  true  "Coupon data"
// @Router       /coupons [post]
// @Success      201  {object}  response.SuccessWithCoupon
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Coupon code already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CouponController) CreateCoupon(ctx *fiber.Ctx) error {
	req := new(validation.CreateCoupon)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	coupon, err := c.CouponService.CreateCoupon(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.SuccessWithCoupon{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Coupon created successfully",
		Coupon:  *coupon,
	})
}

// @Tags         Coupon
// @Summary      Update coupon
// @Description  Update an existing coupon
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path      string                   true  "Coupon ID"
// @Param        coupon  body      validation.UpdateCoupon  true  "Coupon data to update"
// @Router       /coupons/{id} [put]
// @Success      200  {object}  response.SuccessWithCoupon
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Coupon not found"
// @Failure      409  {object}  response.ErrorDetails  "Coupon code already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CouponController) UpdateCoupon(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid coupon ID format",
			Errors:  "Coupon ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateCoupon)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	coupon, err := c.CouponService.UpdateCoupon(ctx, id, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithCoupon{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Coupon updated successfully",
		Coupon:  *coupon,
	})
}

// @Tags         Coupon
// @Summary      Delete coupon
// @Description  Delete a coupon by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Coupon ID"
// @Router       /coupons/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Coupon not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CouponController) DeleteCoupon(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid coupon ID format",
			Errors:  "Coupon ID must be a valid UUID",
		})
	}

	err := c.CouponService.DeleteCoupon(ctx, id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Coupon deleted successfully",
	})
}
