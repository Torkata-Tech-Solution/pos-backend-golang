package validation

// Customer validations
type CreateCustomer struct {
	UserID        string `json:"user_id" validate:"omitempty,uuid"`
	Name          string `json:"name" validate:"required,max=255"`
	Email         string `json:"email" validate:"required,email,max=255"`
	Phone         string `json:"phone" validate:"omitempty,max=50"`
	Address       string `json:"address" validate:"omitempty"`
	LoyaltyPoints int    `json:"loyalty_points" validate:"omitempty,min=0"`
	OutletID      string `json:"outlet_id" validate:"required,uuid"`
}

type UpdateCustomer struct {
	Name          string `json:"name" validate:"omitempty,max=255"`
	Email         string `json:"email" validate:"omitempty,email,max=255"`
	Phone         string `json:"phone" validate:"omitempty,max=50"`
	Address       string `json:"address" validate:"omitempty"`
	LoyaltyPoints int    `json:"loyalty_points" validate:"omitempty,min=0"`
}
