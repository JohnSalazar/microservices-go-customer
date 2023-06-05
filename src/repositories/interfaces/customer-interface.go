package interfaces

import (
	"context"
	"customer/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.Customer, error)
	FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Customer, error)
	Create(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	Delete(ctx context.Context, ID primitive.ObjectID) error
}
