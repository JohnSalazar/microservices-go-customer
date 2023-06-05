package services

import (
	"context"
	"customer/src/models"
	"customer/src/repositories/interfaces"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerService struct {
	customerRepository interfaces.CustomerRepository
}

func NewCustomerService(
	customerRepository interfaces.CustomerRepository,
) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
	}
}

func (service *CustomerService) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	return service.customerRepository.FindByEmail(ctx, email)
}

func (service *CustomerService) FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Customer, error) {
	return service.customerRepository.FindByID(ctx, ID)
}

func (service *CustomerService) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer, err := service.customerRepository.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (service *CustomerService) Update(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer, err := service.customerRepository.Update(ctx, customer)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	return customer, nil
}

func (service *CustomerService) Delete(ctx context.Context, ID primitive.ObjectID) error {
	return service.customerRepository.Delete(ctx, ID)
}
