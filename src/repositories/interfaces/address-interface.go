package interfaces

import (
	"context"
	"customer/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressRepository interface {
	GetByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*models.Address, error)
	FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Address, error)
	FindAddressExists(ctx context.Context, address *models.Address) bool
	Create(ctx context.Context, address *models.Address) error
	Update(ctx context.Context, address *models.Address) (*models.Address, error)
	Delete(ctx context.Context, ID primitive.ObjectID) error
}
