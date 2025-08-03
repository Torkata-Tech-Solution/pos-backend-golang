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

type CustomerService interface {
	GetCustomers(c *fiber.Ctx, params *validation.QueryParams) ([]model.Customer, int64, error)
	GetCustomerByID(c *fiber.Ctx, id string) (*model.Customer, error)
	GetCustomersByOutletID(c *fiber.Ctx, outletID string) ([]model.Customer, error)
	GetCustomerByEmail(c *fiber.Ctx, email string) (*model.Customer, error)
	CreateCustomer(c *fiber.Ctx, req *validation.CreateCustomer) (*model.Customer, error)
	UpdateCustomer(c *fiber.Ctx, id string, req *validation.UpdateCustomer) (*model.Customer, error)
	DeleteCustomer(c *fiber.Ctx, id string) error
	AddLoyaltyPoints(c *fiber.Ctx, id string, points int) error
}

type customerService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewCustomerService(db *gorm.DB, validate *validator.Validate) CustomerService {
	return &customerService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *customerService) GetCustomers(c *fiber.Ctx, params *validation.QueryParams) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("User").Preload("Outlet").Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&customers).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search customers: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&customers).Error; err != nil {
		s.Log.Errorf("Failed to get customers: %+v", err)
		return nil, 0, err
	}

	return customers, totalResults, nil
}

func (s *customerService) GetCustomerByID(c *fiber.Ctx, id string) (*model.Customer, error) {
	customer := new(model.Customer)
	result := s.DB.WithContext(c.Context()).Preload("User").Preload("Outlet").First(customer, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Customer not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get customer by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return customer, nil
}

func (s *customerService) GetCustomersByOutletID(c *fiber.Ctx, outletID string) ([]model.Customer, error) {
	var customers []model.Customer
	result := s.DB.WithContext(c.Context()).Preload("User").Where("outlet_id = ?", outletID).Find(&customers)
	if result.Error != nil {
		s.Log.Errorf("Failed to get customers by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return customers, nil
}

func (s *customerService) GetCustomerByEmail(c *fiber.Ctx, email string) (*model.Customer, error) {
	customer := new(model.Customer)
	result := s.DB.WithContext(c.Context()).Preload("User").Preload("Outlet").Where("email = ?", email).First(customer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Customer not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get customer by email %s: %+v", email, result.Error)
		return nil, result.Error
	}
	return customer, nil
}

func (s *customerService) CreateCustomer(c *fiber.Ctx, req *validation.CreateCustomer) (*model.Customer, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	customer := &model.Customer{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         &req.Phone,
		Address:       &req.Address,
		LoyaltyPoints: req.LoyaltyPoints,
		OutletID:      uuid.MustParse(req.OutletID),
	}

	// Set UserID if provided
	if req.UserID != "" {
		userID := uuid.MustParse(req.UserID)
		customer.UserID = &userID
	}

	result := s.DB.WithContext(c.Context()).Create(customer)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Customer with this email already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create customer: %+v", result.Error)
	}

	return customer, result.Error
}

func (s *customerService) UpdateCustomer(c *fiber.Ctx, id string, req *validation.UpdateCustomer) (*model.Customer, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Name == "" && req.Email == "" && req.Phone == "" && req.Address == "" && req.LoyaltyPoints == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	customer := &model.Customer{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         &req.Phone,
		Address:       &req.Address,
		LoyaltyPoints: req.LoyaltyPoints,
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Customer{}).Where("id = ?", id).Updates(customer)
	if result.Error != nil {
		s.Log.Errorf("Failed to update customer with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Customer not found")
	}

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Customer with this email already exists")
	}

	customer, err := s.GetCustomerByID(c, id)
	if err != nil {
		return nil, err
	}

	return customer, result.Error
}

func (s *customerService) DeleteCustomer(c *fiber.Ctx, id string) error {
	customer := new(model.Customer)

	result := s.DB.WithContext(c.Context()).Delete(customer, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete customer with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Customer not found")
	}

	return nil
}

func (s *customerService) AddLoyaltyPoints(c *fiber.Ctx, id string, points int) error {
	result := s.DB.WithContext(c.Context()).Model(&model.Customer{}).
		Where("id = ?", id).
		Update("loyalty_points", gorm.Expr("loyalty_points + ?", points))

	if result.Error != nil {
		s.Log.Errorf("Failed to add loyalty points to customer %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Customer not found")
	}

	return nil
}
