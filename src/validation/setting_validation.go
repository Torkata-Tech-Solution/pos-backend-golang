package validation

// Setting validations
type CreateSetting struct {
	OutletID string `json:"outlet_id" validate:"required,uuid"`
	Key      string `json:"key" validate:"required,max=100"`
	Value    string `json:"value" validate:"required"`
}

type UpdateSetting struct {
	Value string `json:"value" validate:"omitempty"`
}
