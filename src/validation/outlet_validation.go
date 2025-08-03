package validation

// Outlet validations
type CreateOutlet struct {
	BusinessID string `json:"business_id" validate:"required,uuid"`
	Name       string `json:"name" validate:"required,max=255"`
	Address    string `json:"address" validate:"required,max=255"`
	Phone      string `json:"phone" validate:"omitempty,max=20"`
	Email      string `json:"email" validate:"omitempty,email,max=255"`
}

type UpdateOutlet struct {
	Name    string `json:"name" validate:"omitempty,max=255"`
	Address string `json:"address" validate:"omitempty,max=255"`
	Phone   string `json:"phone" validate:"omitempty,max=20"`
	Email   string `json:"email" validate:"omitempty,email,max=255"`
}

// Outlet Staff validations
type CreateOutletStaff struct {
	OutletID string `json:"outlet_id" validate:"required,uuid"`
	Name     string `json:"name" validate:"required,max=255"`
	Username string `json:"username" validate:"omitempty,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	Role     string `json:"role" validate:"required,max=255"`
}

type UpdateOutletStaff struct {
	Name     string `json:"name" validate:"omitempty,max=255"`
	Username string `json:"username" validate:"omitempty,max=255"`
	Password string `json:"password" validate:"omitempty,min=8,max=255"`
	Role     string `json:"role" validate:"omitempty,max=255"`
}

type ChangePasswordOutletStaff struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=255"`
}
