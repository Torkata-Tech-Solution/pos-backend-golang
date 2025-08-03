package validation

import "time"

// Sale validations
type CreateSale struct {
	OutletID        string    `json:"outlet_id" validate:"required,uuid"`
	OutletStaffID   string    `json:"outlet_staff_id" validate:"required,uuid"`
	CustomerID      string    `json:"customer_id" validate:"omitempty,uuid"`
	PaymentMethodID string    `json:"payment_method_id" validate:"omitempty,uuid"`
	TableID         string    `json:"table_id" validate:"required,uuid"`
	InvoiceNumber   string    `json:"invoice_number" validate:"required,max=50"`
	Total           float64   `json:"total" validate:"required,min=0"`
	Discount        float64   `json:"discount" validate:"omitempty,min=0"`
	Tax             float64   `json:"tax" validate:"omitempty,min=0"`
	GrandTotal      float64   `json:"grand_total" validate:"required,min=0"`
	Status          string    `json:"status" validate:"required,oneof=paid unpaid void hold"`
	SaleDate        time.Time `json:"sale_date" validate:"omitempty"`
	Note            string    `json:"note" validate:"omitempty"`
}

type UpdateSale struct {
	CustomerID      string  `json:"customer_id" validate:"omitempty,uuid"`
	PaymentMethodID string  `json:"payment_method_id" validate:"omitempty,uuid"`
	TableID         string  `json:"table_id" validate:"omitempty,uuid"`
	Total           float64 `json:"total" validate:"omitempty,min=0"`
	Discount        float64 `json:"discount" validate:"omitempty,min=0"`
	Tax             float64 `json:"tax" validate:"omitempty,min=0"`
	GrandTotal      float64 `json:"grand_total" validate:"omitempty,min=0"`
	Status          string  `json:"status" validate:"omitempty,oneof=paid unpaid void hold"`
	Note            string  `json:"note" validate:"omitempty"`
}

// Sale Item validations
type CreateSaleItem struct {
	SaleID    string  `json:"sale_id" validate:"required,uuid"`
	ProductID string  `json:"product_id" validate:"required,uuid"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=0"`
	Discount  float64 `json:"discount" validate:"omitempty,min=0"`
	Total     float64 `json:"total" validate:"required,min=0"`
}

type UpdateSaleItem struct {
	Quantity int     `json:"quantity" validate:"omitempty,min=1"`
	Price    float64 `json:"price" validate:"omitempty,min=0"`
	Discount float64 `json:"discount" validate:"omitempty,min=0"`
	Total    float64 `json:"total" validate:"omitempty,min=0"`
}

type SaleFilters struct {
	OutletID      string `json:"outlet_id" validate:"omitempty,uuid"`
	CustomerID    string `json:"customer_id" validate:"omitempty,uuid"`
	PaymentMethod string `json:"payment_method" validate:"omitempty,uuid"`
	Status        string `json:"status" validate:"omitempty,oneof=paid unpaid void hold"`
	DateFrom      string `json:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DateTo        string `json:"date_to" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateSaleStatus struct {
	Status string `json:"status" validate:"required,oneof=paid unpaid void hold"`
}

type SalesReportFilters struct {
	OutletID string `json:"outlet_id" validate:"omitempty,uuid"`
	DateFrom string `json:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DateTo   string `json:"date_to" validate:"omitempty,datetime=2006-01-02"`
	GroupBy  string `json:"group_by" validate:"omitempty,oneof=day month year"`
}
