package validation

// Product Category validations
type CreateProductCategory struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty"`
	BusinessID  string `json:"business_id" validate:"required,uuid"`
}

type UpdateProductCategory struct {
	Name        string `json:"name" validate:"omitempty,max=255"`
	Description string `json:"description" validate:"omitempty"`
}

// Product validations
type CreateProduct struct {
	Image       string  `json:"image" validate:"omitempty,url,max=255"`
	Name        string  `json:"name" validate:"required,max=255"`
	Description string  `json:"description" validate:"omitempty"`
	Price       float64 `json:"price" validate:"required,min=0"`
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	BusinessID  string  `json:"business_id" validate:"required,uuid"`
}

type UpdateProduct struct {
	Image       string  `json:"image" validate:"omitempty,url,max=255"`
	Name        string  `json:"name" validate:"omitempty,max=255"`
	Description string  `json:"description" validate:"omitempty"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
	CategoryID  string  `json:"category_id" validate:"omitempty,uuid"`
}
