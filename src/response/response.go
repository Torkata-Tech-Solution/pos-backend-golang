package response

import "app/src/model"

// Common response types
type Common struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SuccessWithUser struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
}

type SuccessWithTokens struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
	Tokens  Tokens     `json:"tokens"`
}

// Paginated response types for each model
type SuccessWithPaginatedUsers struct {
	Code         int          `json:"code"`
	Status       string       `json:"status"`
	Message      string       `json:"message"`
	Results      []model.User `json:"results"`
	Page         int          `json:"page"`
	Limit        int          `json:"limit"`
	TotalPages   int64        `json:"total_pages"`
	TotalResults int64        `json:"total_results"`
}

type SuccessWithPaginatedBusinesses struct {
	Code         int              `json:"code"`
	Status       string           `json:"status"`
	Message      string           `json:"message"`
	Results      []model.Business `json:"results"`
	Page         int              `json:"page"`
	Limit        int              `json:"limit"`
	TotalPages   int64            `json:"total_pages"`
	TotalResults int64            `json:"total_results"`
}

type SuccessWithPaginatedOutlets struct {
	Code         int            `json:"code"`
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	Results      []model.Outlet `json:"results"`
	Page         int            `json:"page"`
	Limit        int            `json:"limit"`
	TotalPages   int64          `json:"total_pages"`
	TotalResults int64          `json:"total_results"`
}

type SuccessWithPaginatedProducts struct {
	Code         int             `json:"code"`
	Status       string          `json:"status"`
	Message      string          `json:"message"`
	Results      []model.Product `json:"results"`
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
	TotalPages   int64           `json:"total_pages"`
	TotalResults int64           `json:"total_results"`
}

type SuccessWithPaginatedProductCategories struct {
	Code         int                     `json:"code"`
	Status       string                  `json:"status"`
	Message      string                  `json:"message"`
	Results      []model.ProductCategory `json:"results"`
	Page         int                     `json:"page"`
	Limit        int                     `json:"limit"`
	TotalPages   int64                   `json:"total_pages"`
	TotalResults int64                   `json:"total_results"`
}

type SuccessWithPaginatedCustomers struct {
	Code         int              `json:"code"`
	Status       string           `json:"status"`
	Message      string           `json:"message"`
	Results      []model.Customer `json:"results"`
	Page         int              `json:"page"`
	Limit        int              `json:"limit"`
	TotalPages   int64            `json:"total_pages"`
	TotalResults int64            `json:"total_results"`
}

type SuccessWithPaginatedTables struct {
	Code         int           `json:"code"`
	Status       string        `json:"status"`
	Message      string        `json:"message"`
	Results      []model.Table `json:"results"`
	Page         int           `json:"page"`
	Limit        int           `json:"limit"`
	TotalPages   int64         `json:"total_pages"`
	TotalResults int64         `json:"total_results"`
}

type SuccessWithPaginatedPaymentMethods struct {
	Code         int                   `json:"code"`
	Status       string                `json:"status"`
	Message      string                `json:"message"`
	Results      []model.PaymentMethod `json:"results"`
	Page         int                   `json:"page"`
	Limit        int                   `json:"limit"`
	TotalPages   int64                 `json:"total_pages"`
	TotalResults int64                 `json:"total_results"`
}

type SuccessWithPaginatedPrinters struct {
	Code         int             `json:"code"`
	Status       string          `json:"status"`
	Message      string          `json:"message"`
	Results      []model.Printer `json:"results"`
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
	TotalPages   int64           `json:"total_pages"`
	TotalResults int64           `json:"total_results"`
}

type SuccessWithPaginatedSettings struct {
	Code         int             `json:"code"`
	Status       string          `json:"status"`
	Message      string          `json:"message"`
	Results      []model.Setting `json:"results"`
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
	TotalPages   int64           `json:"total_pages"`
	TotalResults int64           `json:"total_results"`
}

type SuccessWithPaginatedCoupons struct {
	Code         int            `json:"code"`
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	Results      []model.Coupon `json:"results"`
	Page         int            `json:"page"`
	Limit        int            `json:"limit"`
	TotalPages   int64          `json:"total_pages"`
	TotalResults int64          `json:"total_results"`
}

type SuccessWithPaginatedSales struct {
	Code         int          `json:"code"`
	Status       string       `json:"status"`
	Message      string       `json:"message"`
	Results      []model.Sale `json:"results"`
	Page         int          `json:"page"`
	Limit        int          `json:"limit"`
	TotalPages   int64        `json:"total_pages"`
	TotalResults int64        `json:"total_results"`
}

type SuccessWithPaginatedOutletStaff struct {
	Code         int                 `json:"code"`
	Status       string              `json:"status"`
	Message      string              `json:"message"`
	Results      []model.OutletStaff `json:"results"`
	Page         int                 `json:"page"`
	Limit        int                 `json:"limit"`
	TotalPages   int64               `json:"total_pages"`
	TotalResults int64               `json:"total_results"`
}

type SuccessWithPaginatedBusinessUsers struct {
	Code         int                  `json:"code"`
	Status       string               `json:"status"`
	Message      string               `json:"message"`
	Results      []model.BusinessUser `json:"results"`
	Page         int                  `json:"page"`
	Limit        int                  `json:"limit"`
	TotalPages   int64                `json:"total_pages"`
	TotalResults int64                `json:"total_results"`
}

type ErrorDetails struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

type SuccessWithBusiness struct {
	Code     int            `json:"code"`
	Status   string         `json:"status"`
	Message  string         `json:"message"`
	Business model.Business `json:"business"`
}

type SuccessWithBusinessUser struct {
	Code         int                `json:"code"`
	Status       string             `json:"status"`
	Message      string             `json:"message"`
	BusinessUser model.BusinessUser `json:"business_user"`
}

type SuccessWithCoupon struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Coupon  model.Coupon `json:"coupon"`
}

type SuccessWithCustomer struct {
	Code     int            `json:"code"`
	Status   string         `json:"status"`
	Message  string         `json:"message"`
	Customer model.Customer `json:"customer"`
}

type SuccessWithOutlet struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Outlet  model.Outlet `json:"outlet"`
}

type SuccessWithOutletStaff struct {
	Code        int               `json:"code"`
	Status      string            `json:"status"`
	Message     string            `json:"message"`
	OutletStaff model.OutletStaff `json:"outlet_staff"`
}

type SuccessWithPaymentMethod struct {
	Code          int                 `json:"code"`
	Status        string              `json:"status"`
	Message       string              `json:"message"`
	PaymentMethod model.PaymentMethod `json:"payment_method"`
}

type SuccessWithProduct struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Product model.Product `json:"product"`
}

type SuccessWithProductCategory struct {
	Code            int                   `json:"code"`
	Status          string                `json:"status"`
	Message         string                `json:"message"`
	ProductCategory model.ProductCategory `json:"product_category"`
}

type SuccessWithPrinter struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Printer model.Printer `json:"printer"`
}

type SuccessWithSale struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Sale    model.Sale `json:"sale"`
}

type SuccessWithSetting struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Setting model.Setting `json:"setting"`
}

type SuccessWithTable struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Table   model.Table `json:"table"`
}

type SuccessWithSalesReport struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Report  []model.Sale `json:"report"`
}
