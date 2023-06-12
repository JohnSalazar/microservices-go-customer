package validators

import (
	"customer/src/dtos"

	common_validator "github.com/JohnSalazar/microservices-go-common/validators"
)

type addAddress struct {
	Street   string `from:"street" json:"street" validate:"required,max=300"`
	City     string `from:"city" json:"city" validate:"required,max=150"`
	Province string `from:"province" json:"province" validate:"required,max=150"`
	Code     string `from:"code" json:"code" validate:"required,max=20"`
	Type     string `from:"type" json:"type" validate:"required,oneof='home' 'billing'"`
}

type updateAddress struct {
	ID       string `from:"id" json:"id" validate:"required"`
	Street   string `from:"street" json:"street" validate:"required,max=300"`
	City     string `from:"city" json:"city" validate:"required,max=150"`
	Province string `from:"province" json:"province" validate:"required,max=150"`
	Code     string `from:"code" json:"code" validate:"required,max=20"`
	Type     string `from:"type" json:"type" validate:"required,oneof='home' 'billing'"`
}

func ValidateAddAddress(fields *dtos.AddAddress) interface{} {
	addAddress := addAddress{
		Street:   fields.Street,
		City:     fields.City,
		Province: fields.Province,
		Code:     fields.Code,
		Type:     fields.Type,
	}

	err := common_validator.Validate(addAddress)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUpdateAddress(fields *dtos.UpdateAddress) interface{} {
	updateAddress := updateAddress{
		ID:       fields.ID.Hex(),
		Street:   fields.Street,
		City:     fields.City,
		Province: fields.Province,
		Code:     fields.Code,
		Type:     fields.Type,
	}

	err := common_validator.Validate(updateAddress)
	if err != nil {
		return err
	}

	return nil
}
