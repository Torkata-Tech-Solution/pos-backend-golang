package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SaleService interface {
	GetSales(c *fiber.Ctx, params *validation.QueryParams, filters *validation.SaleFilters) ([]model.Sale, int64, error)
	GetSaleByID(c *fiber.Ctx, id string) (*model.Sale, error)
	GetSaleByInvoiceNumber(c *fiber.Ctx, invoiceNumber string) (*model.Sale, error)
	GetSalesByOutletID(c *fiber.Ctx, outletID string) ([]model.Sale, error)
	GetSalesByDateRange(c *fiber.Ctx, outletID string, startDate, endDate time.Time) ([]model.Sale, error)
	CreateSale(c *fiber.Ctx, req *validation.CreateSale) (*model.Sale, error)
	UpdateSale(c *fiber.Ctx, id string, req *validation.UpdateSale) (*model.Sale, error)
	DeleteSale(c *fiber.Ctx, id string) error
	UpdateSaleStatus(c *fiber.Ctx, id string, req *validation.UpdateSaleStatus) (*model.Sale, error)
	GetSalesReport(c *fiber.Ctx, filters *validation.SalesReportFilters) ([]model.Sale, error)
}

type SaleItemService interface {
	GetSaleItems(c *fiber.Ctx, saleID string) ([]model.SaleItem, error)
	GetSaleItemByID(c *fiber.Ctx, id string) (*model.SaleItem, error)
	CreateSaleItem(c *fiber.Ctx, req *validation.CreateSaleItem) (*model.SaleItem, error)
	UpdateSaleItem(c *fiber.Ctx, id string, req *validation.UpdateSaleItem) (*model.SaleItem, error)
	DeleteSaleItem(c *fiber.Ctx, id string) error
}

// Sale Service Implementation
type saleService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewSaleService(db *gorm.DB, validate *validator.Validate) SaleService {
	return &saleService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *saleService) GetSales(c *fiber.Ctx, params *validation.QueryParams, filters *validation.SaleFilters) ([]model.Sale, int64, error) {
	var sales []model.Sale
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).
		Preload("Outlet").
		Preload("OutletStaff").
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Table").
		Preload("SaleItems").
		Order("created_at desc")

	if search := params.Search; search != "" {
		query = query.Where("invoice_number LIKE ? OR note LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&sales).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search sales: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&sales).Error; err != nil {
		s.Log.Errorf("Failed to get sales: %+v", err)
		return nil, 0, err
	}

	return sales, totalResults, nil
}

func (s *saleService) GetSaleByID(c *fiber.Ctx, id string) (*model.Sale, error) {
	sale := new(model.Sale)
	result := s.DB.WithContext(c.Context()).
		Preload("Outlet").
		Preload("OutletStaff").
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Table").
		Preload("SaleItems").
		Preload("SaleItems.Product").
		First(sale, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sale not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get sale by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return sale, nil
}

func (s *saleService) GetSaleByInvoiceNumber(c *fiber.Ctx, invoiceNumber string) (*model.Sale, error) {
	sale := new(model.Sale)
	result := s.DB.WithContext(c.Context()).
		Preload("Outlet").
		Preload("OutletStaff").
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Table").
		Preload("SaleItems").
		Preload("SaleItems.Product").
		Where("invoice_number = ?", invoiceNumber).
		First(sale)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sale not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get sale by invoice number %s: %+v", invoiceNumber, result.Error)
		return nil, result.Error
	}
	return sale, nil
}

func (s *saleService) GetSalesByOutletID(c *fiber.Ctx, outletID string) ([]model.Sale, error) {
	var sales []model.Sale
	result := s.DB.WithContext(c.Context()).
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Table").
		Where("outlet_id = ?", outletID).
		Order("created_at desc").
		Find(&sales)

	if result.Error != nil {
		s.Log.Errorf("Failed to get sales by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return sales, nil
}

func (s *saleService) GetSalesByDateRange(c *fiber.Ctx, outletID string, startDate, endDate time.Time) ([]model.Sale, error) {
	var sales []model.Sale
	result := s.DB.WithContext(c.Context()).
		Preload("Customer").
		Preload("PaymentMethod").
		Preload("Table").
		Where("outlet_id = ? AND sale_date BETWEEN ? AND ?", outletID, startDate, endDate).
		Order("sale_date desc").
		Find(&sales)

	if result.Error != nil {
		s.Log.Errorf("Failed to get sales by date range: %+v", result.Error)
		return nil, result.Error
	}
	return sales, nil
}

func (s *saleService) CreateSale(c *fiber.Ctx, req *validation.CreateSale) (*model.Sale, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	saleDate := time.Now()
	if !req.SaleDate.IsZero() {
		saleDate = req.SaleDate
	}

	sale := &model.Sale{
		OutletID:      uuid.MustParse(req.OutletID),
		OutletStaffID: uuid.MustParse(req.OutletStaffID),
		TableID:       uuid.MustParse(req.TableID),
		InvoiceNumber: req.InvoiceNumber,
		Total:         req.Total,
		Discount:      req.Discount,
		Tax:           req.Tax,
		GrandTotal:    req.GrandTotal,
		Status:        req.Status,
		SaleDate:      saleDate,
		Note:          &req.Note,
	}

	// Set optional fields if provided
	if req.CustomerID != "" {
		customerID := uuid.MustParse(req.CustomerID)
		sale.CustomerID = &customerID
	}
	if req.PaymentMethodID != "" {
		paymentMethodID := uuid.MustParse(req.PaymentMethodID)
		sale.PaymentMethodID = &paymentMethodID
	}

	result := s.DB.WithContext(c.Context()).Create(sale)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Sale with this invoice number already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create sale: %+v", result.Error)
	}

	return sale, result.Error
}

func (s *saleService) UpdateSale(c *fiber.Ctx, id string, req *validation.UpdateSale) (*model.Sale, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	updateFields := make(map[string]interface{})

	if req.CustomerID != "" {
		updateFields["customer_id"] = req.CustomerID
	}
	if req.PaymentMethodID != "" {
		updateFields["payment_method_id"] = req.PaymentMethodID
	}
	if req.TableID != "" {
		updateFields["table_id"] = req.TableID
	}
	if req.Total > 0 {
		updateFields["total"] = req.Total
	}
	if req.Discount >= 0 {
		updateFields["discount"] = req.Discount
	}
	if req.Tax >= 0 {
		updateFields["tax"] = req.Tax
	}
	if req.GrandTotal > 0 {
		updateFields["grand_total"] = req.GrandTotal
	}
	if req.Status != "" {
		updateFields["status"] = req.Status
	}
	if req.Note != "" {
		updateFields["note"] = req.Note
	}

	if len(updateFields) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Sale{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		s.Log.Errorf("Failed to update sale with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sale not found")
	}

	sale, err := s.GetSaleByID(c, id)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *saleService) DeleteSale(c *fiber.Ctx, id string) error {
	sale := new(model.Sale)

	result := s.DB.WithContext(c.Context()).Delete(sale, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete sale with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Sale not found")
	}

	return nil
}

func (s *saleService) UpdateSaleStatus(c *fiber.Ctx, id string, req *validation.UpdateSaleStatus) (*model.Sale, error) {
	result := s.DB.WithContext(c.Context()).Model(&model.Sale{}).
		Where("id = ?", id).
		Update("status", req.Status)

	if result.Error != nil {
		s.Log.Errorf("Failed to update sale status for ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sale not found")
	}

	sale, err := s.GetSaleByID(c, id)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *saleService) GetSalesReport(c *fiber.Ctx, filters *validation.SalesReportFilters) ([]model.Sale, error) {
	var sales []model.Sale

	// Apply filters to the query
	query := s.DB.WithContext(c.Context()).Model(&model.Sale{})

	if filters.OutletID != "" {
		query = query.Where("outlet_id = ?", filters.OutletID)
	}
	if filters.DateFrom != "" {
		query = query.Where("sale_date >= ?", filters.DateFrom)
	}
	if filters.DateTo != "" {
		query = query.Where("sale_date <= ?", filters.DateTo)
	}
	// if filters.GroupBy != "" {
	// 	query = query.Group("DATE_TRUNC(?, sale_date)", filters.GroupBy)
	// }

	result := query.Find(&sales)
	if result.Error != nil {
		s.Log.Errorf("Failed to get sales report: %+v", result.Error)
		return nil, result.Error
	}

	return sales, nil
}

// Sale Item Service Implementation
type saleItemService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewSaleItemService(db *gorm.DB, validate *validator.Validate) SaleItemService {
	return &saleItemService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *saleItemService) GetSaleItems(c *fiber.Ctx, saleID string) ([]model.SaleItem, error) {
	var saleItems []model.SaleItem
	result := s.DB.WithContext(c.Context()).
		Preload("Product").
		Where("sale_id = ?", saleID).
		Find(&saleItems)

	if result.Error != nil {
		s.Log.Errorf("Failed to get sale items for sale ID %s: %+v", saleID, result.Error)
		return nil, result.Error
	}
	return saleItems, nil
}

func (s *saleItemService) GetSaleItemByID(c *fiber.Ctx, id string) (*model.SaleItem, error) {
	saleItem := new(model.SaleItem)
	result := s.DB.WithContext(c.Context()).
		Preload("Product").
		Preload("Sale").
		First(saleItem, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sale item not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get sale item by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return saleItem, nil
}

func (s *saleItemService) CreateSaleItem(c *fiber.Ctx, req *validation.CreateSaleItem) (*model.SaleItem, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	saleItem := &model.SaleItem{
		SaleID:    uuid.MustParse(req.SaleID),
		ProductID: uuid.MustParse(req.ProductID),
		Quantity:  req.Quantity,
		Price:     req.Price,
		Discount:  req.Discount,
		Total:     req.Total,
	}

	result := s.DB.WithContext(c.Context()).Create(saleItem)
	if result.Error != nil {
		s.Log.Errorf("Failed to create sale item: %+v", result.Error)
	}

	return saleItem, result.Error
}

func (s *saleItemService) UpdateSaleItem(c *fiber.Ctx, id string, req *validation.UpdateSaleItem) (*model.SaleItem, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	updateFields := make(map[string]interface{})

	if req.Quantity > 0 {
		updateFields["quantity"] = req.Quantity
	}
	if req.Price > 0 {
		updateFields["price"] = req.Price
	}
	if req.Discount >= 0 {
		updateFields["discount"] = req.Discount
	}
	if req.Total > 0 {
		updateFields["total"] = req.Total
	}

	if len(updateFields) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	result := s.DB.WithContext(c.Context()).Model(&model.SaleItem{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		s.Log.Errorf("Failed to update sale item with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sale item not found")
	}

	saleItem, err := s.GetSaleItemByID(c, id)
	if err != nil {
		return nil, err
	}

	return saleItem, nil
}

func (s *saleItemService) DeleteSaleItem(c *fiber.Ctx, id string) error {
	saleItem := new(model.SaleItem)

	result := s.DB.WithContext(c.Context()).Delete(saleItem, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete sale item with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Sale item not found")
	}

	return nil
}
