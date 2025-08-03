package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PrinterService interface {
	GetPrinters(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Printer, int64, error)
	GetPrinterByID(c *fiber.Ctx, id string) (*model.Printer, error)
	GetPrintersByOutletID(c *fiber.Ctx, outletID string) ([]model.Printer, error)
	GetDefaultPrinterByOutletID(c *fiber.Ctx, outletID string) (*model.Printer, error)
	CreatePrinter(c *fiber.Ctx, req *validation.CreatePrinter) (*model.Printer, error)
	UpdatePrinter(c *fiber.Ctx, id string, req *validation.UpdatePrinter) (*model.Printer, error)
	DeletePrinter(c *fiber.Ctx, id string) error
}

type printerService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewPrinterService(db *gorm.DB, validate *validator.Validate) PrinterService {
	return &printerService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *printerService) GetPrinters(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Printer, int64, error) {
	var printers []model.Printer
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at desc")

	if search := params.Search; search != "" {
		query = query.Where("name ILIKE ? OR connection_type ILIKE ? OR mac_address ILIKE ? OR ip_address ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	result := query.Model(&model.Printer{}).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to count printers: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&printers).Error; err != nil {
		s.Log.Errorf("Failed to get printers: %+v", err)
		return nil, 0, err
	}

	return printers, totalResults, nil
}

func (s *printerService) GetPrinterByID(c *fiber.Ctx, id string) (*model.Printer, error) {
	printer := new(model.Printer)
	result := s.DB.WithContext(c.Context()).Preload("Outlet").First(printer, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Printer not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get printer by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return printer, nil
}

func (s *printerService) GetPrintersByOutletID(c *fiber.Ctx, outletID string) ([]model.Printer, error) {
	var printers []model.Printer
	result := s.DB.WithContext(c.Context()).Where("outlet_id = ?", outletID).Find(&printers)
	if result.Error != nil {
		s.Log.Errorf("Failed to get printers by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return printers, nil
}

func (s *printerService) GetDefaultPrinterByOutletID(c *fiber.Ctx, outletID string) (*model.Printer, error) {
	printer := new(model.Printer)
	result := s.DB.WithContext(c.Context()).Where("outlet_id = ? AND default_printer = true", outletID).First(printer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Default printer not found for this outlet")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get default printer by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return printer, nil
}

func (s *printerService) CreatePrinter(c *fiber.Ctx, req *validation.CreatePrinter) (*model.Printer, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Check if outlet exists
	var outlet model.Outlet
	if err := s.DB.WithContext(c.Context()).First(&outlet, "id = ?", req.OutletID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Outlet not found")
		}
		return nil, err
	}

	// If this is to be the default printer, unset all other default printers for this outlet
	if req.DefaultPrinter {
		if err := s.DB.WithContext(c.Context()).Model(&model.Printer{}).
			Where("outlet_id = ?", req.OutletID).
			Update("default_printer", false).Error; err != nil {
			s.Log.Errorf("Failed to unset default printers: %+v", err)
			return nil, err
		}
	}

	printer := &model.Printer{
		OutletID:       uuid.MustParse(req.OutletID),
		Name:           req.Name,
		ConnectionType: req.ConnectionType,
		DefaultPrinter: req.DefaultPrinter,
	}

	// Set optional fields if provided
	if req.MacAddress != "" {
		printer.MacAddress = &req.MacAddress
	}
	if req.IPAddress != "" {
		printer.IPAddress = &req.IPAddress
	}
	if req.PaperWidth > 0 {
		printer.PaperWidth = &req.PaperWidth
	}

	result := s.DB.WithContext(c.Context()).Create(printer)
	if result.Error != nil {
		s.Log.Errorf("Failed to create printer: %+v", result.Error)
		return nil, result.Error
	}

	return printer, nil
}

func (s *printerService) UpdatePrinter(c *fiber.Ctx, id string, req *validation.UpdatePrinter) (*model.Printer, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	printer := new(model.Printer)
	result := s.DB.WithContext(c.Context()).First(printer, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Printer not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get printer by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	// If setting as default printer, unset all other default printers for this outlet
	if req.DefaultPrinter {
		if err := s.DB.WithContext(c.Context()).Model(&model.Printer{}).
			Where("outlet_id = ? AND id != ?", printer.OutletID, id).
			Update("default_printer", false).Error; err != nil {
			s.Log.Errorf("Failed to unset default printers: %+v", err)
			return nil, err
		}
	}

	// Update fields
	if req.Name != "" {
		printer.Name = req.Name
	}
	if req.ConnectionType != "" {
		printer.ConnectionType = req.ConnectionType
	}
	if req.MacAddress != "" {
		printer.MacAddress = &req.MacAddress
	}
	if req.IPAddress != "" {
		printer.IPAddress = &req.IPAddress
	}
	if req.PaperWidth > 0 {
		printer.PaperWidth = &req.PaperWidth
	}
	printer.DefaultPrinter = req.DefaultPrinter

	result = s.DB.WithContext(c.Context()).Save(printer)
	if result.Error != nil {
		s.Log.Errorf("Failed to update printer: %+v", result.Error)
		return nil, result.Error
	}

	return printer, nil
}

func (s *printerService) DeletePrinter(c *fiber.Ctx, id string) error {
	result := s.DB.WithContext(c.Context()).Delete(&model.Printer{}, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete printer: %+v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Printer not found")
	}
	return nil
}
