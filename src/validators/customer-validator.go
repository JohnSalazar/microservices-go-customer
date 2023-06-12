package validators

import (
	"customer/src/dtos"

	common_validator "github.com/JohnSalazar/microservices-go-common/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type addCustomer struct {
	ID        primitive.ObjectID `from:"id" json:"id" validate:"required"`
	Email     string             `from:"email" json:"email" validate:"required,email"`
	FirstName string             `from:"firstname" json:"firstname" validate:"max=150"`
	LastName  string             `from:"lastname" json:"lastname" validate:"max=150"`
	Phone     string             `from:"phone" json:"phone" validate:"max=30"`
}

type updateCustomer struct {
	ID        primitive.ObjectID `from:"id" json:"id" validate:"required"`
	Email     string             `from:"email" json:"email" validate:"required,email"`
	FirstName string             `from:"firstname" json:"firstname" validate:"max=150"`
	LastName  string             `from:"lastname" json:"lastname" validate:"max=150"`
	Phone     string             `from:"phone" json:"phone" validate:"max=30"`
}

func ValidateAddCustomer(id primitive.ObjectID, fields *dtos.AddCustomer) interface{} {
	addCustomer := addCustomer{
		ID:        id,
		Email:     fields.Email,
		FirstName: fields.FirstName,
		LastName:  fields.LastName,
		Phone:     fields.Phone,
	}

	err := common_validator.Validate(addCustomer)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUpdateCustomer(id primitive.ObjectID, fields *dtos.UpdateCustomer) interface{} {
	updateCustomer := updateCustomer{
		ID:        id,
		Email:     fields.Email,
		FirstName: fields.FirstName,
		LastName:  fields.LastName,
		Phone:     fields.Phone,
	}

	err := common_validator.Validate(updateCustomer)
	if err != nil {
		return err
	}

	return nil
}
