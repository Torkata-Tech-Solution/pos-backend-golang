package validation

// Printer validations
type CreatePrinter struct {
	OutletID       string `json:"outlet_id" validate:"required,uuid"`
	Name           string `json:"name" validate:"required,max=255"`
	ConnectionType string `json:"connection_type" validate:"required,max=50"`
	MacAddress     string `json:"mac_address" validate:"omitempty,max=50"`
	IPAddress      string `json:"ip_address" validate:"omitempty,ip,max=50"`
	PaperWidth     int    `json:"paper_width" validate:"omitempty,min=1"`
	DefaultPrinter bool   `json:"default_printer" validate:"omitempty"`
}

type UpdatePrinter struct {
	Name           string `json:"name" validate:"omitempty,max=255"`
	ConnectionType string `json:"connection_type" validate:"omitempty,max=50"`
	MacAddress     string `json:"mac_address" validate:"omitempty,max=50"`
	IPAddress      string `json:"ip_address" validate:"omitempty,ip,max=50"`
	PaperWidth     int    `json:"paper_width" validate:"omitempty,min=1"`
	DefaultPrinter bool   `json:"default_printer" validate:"omitempty"`
}
