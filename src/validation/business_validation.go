package validation

// Business validations
type CreateBusiness struct {
	Domain  string `json:"domain" validate:"required,max=255"`
	Name    string `json:"name" validate:"required,max=255"`
	Address string `json:"address" validate:"required,max=255"`
	Phone   string `json:"phone" validate:"omitempty,max=20"`
	Email   string `json:"email" validate:"omitempty,email,max=255"`
	Website string `json:"website" validate:"omitempty,url,max=255"`
	Logo    string `json:"logo" validate:"omitempty,url,max=255"`
}

type UpdateBusiness struct {
	Name    string `json:"name" validate:"omitempty,max=255"`
	Address string `json:"address" validate:"omitempty,max=255"`
	Phone   string `json:"phone" validate:"omitempty,max=20"`
	Email   string `json:"email" validate:"omitempty,email,max=255"`
	Website string `json:"website" validate:"omitempty,url,max=255"`
	Logo    string `json:"logo" validate:"omitempty,url,max=255"`
}

// Business User validations
type CreateBusinessUser struct {
	BusinessID string `json:"business_id" validate:"required,uuid"`
	UserID     string `json:"user_id" validate:"required,uuid"`
	Role       string `json:"role" validate:"required,max=255"`
}

type UpdateBusinessUser struct {
	Role string `json:"role" validate:"omitempty,max=255"`
}

// Common query validation
type QueryParams struct {
	Page   int    `json:"page" validate:"omitempty,min=1"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Search string `json:"search" validate:"omitempty"`
}
