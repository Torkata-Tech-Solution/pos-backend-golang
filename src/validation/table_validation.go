package validation

// Table validations
type CreateTable struct {
	OutletID string `json:"outlet_id" validate:"required,uuid"`
	Name     string `json:"name" validate:"required,max=255"`
	Location string `json:"location" validate:"omitempty,max=255"`
	Status   string `json:"status" validate:"omitempty,oneof=available occupied reserved"`
	Capacity int    `json:"capacity" validate:"required,min=1"`
}

type UpdateTable struct {
	Name     string `json:"name" validate:"omitempty,max=255"`
	Location string `json:"location" validate:"omitempty,max=255"`
	Status   string `json:"status" validate:"omitempty,oneof=available occupied reserved"`
	Capacity int    `json:"capacity" validate:"omitempty,min=1"`
}
