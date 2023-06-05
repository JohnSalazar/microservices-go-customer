package dtos

type AddAddress struct {
	Street   string `json:"street,omitempty"`
	City     string `json:"city,omitempty"`
	Province string `json:"province,omitempty"`
	Code     string `json:"code,omitempty"`
	Type     string `json:"type,omitempty"`
}
