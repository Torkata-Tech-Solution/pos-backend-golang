package validation

// Payment Method validations
type CreatePaymentMethod struct {
	OutletID string `json:"outlet_id" validate:"required,uuid"`
	Name     string `json:"name" validate:"required,max=255"`
	Type     string `json:"type" validate:"required,max=255"`
}

type UpdatePaymentMethod struct {
	Name string `json:"name" validate:"omitempty,max=255"`
	Type string `json:"type" validate:"omitempty,max=255"`
}
