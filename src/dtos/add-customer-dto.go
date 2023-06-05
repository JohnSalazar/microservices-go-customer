package dtos

type AddCustomer struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Phone     string `json:"phone,omitempty"`
}
