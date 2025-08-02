package validation

type CreateBusiness struct {
	Domain  string `json:"domain" validate:"required,unique"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Website string `json:"website" validate:"omitempty,url"`
	Logo    string `json:"logo" validate:"omitempty,url"`
}

type UpdateBusiness struct {
	Name    string `json:"name" validate:"omitempty"`
	Address string `json:"address" validate:"omitempty"`
	Phone   string `json:"phone" validate:"omitempty"`
	Email   string `json:"email" validate:"omitempty,email"`
	Website string `json:"website" validate:"omitempty,url"`
	Logo    string `json:"logo" validate:"omitempty,url"`
}

type QueryBusiness struct {
	Page   int    `json:"page" validate:"omitempty,number,min=1"`
	Limit  int    `json:"limit" validate:"omitempty,number,min=1"`
	Search string `json:"search" validate:"omitempty"`
}
