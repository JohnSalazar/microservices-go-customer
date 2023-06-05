package services

import (
	"context"
	"customer/src/models"
	"customer/src/repositories/interfaces"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressService struct {
	addressRepository interfaces.AddressRepository
}

func NewAddressService(
	addressRepository interfaces.AddressRepository,
) *AddressService {
	return &AddressService{
		addressRepository: addressRepository,
	}
}

func (service *AddressService) GetByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*models.Address, error) {
	return service.addressRepository.GetByCustomerID(ctx, customerID)
}

func (service *AddressService) FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Address, error) {
	return service.addressRepository.FindByID(ctx, ID)
}

func (service *AddressService) FindAddressExists(ctx context.Context, address *models.Address) bool {
	return service.addressRepository.FindAddressExists(ctx, address)
}

func (service *AddressService) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	address.ID = primitive.NewObjectID()
	address.CreatedAt = time.Now().UTC()

	err := service.addressRepository.Create(ctx, address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (service *AddressService) Update(ctx context.Context, address *models.Address) (*models.Address, error) {
	_address, err := service.addressRepository.Update(ctx, address)
	if err != nil {
		return nil, errors.New("address not found")
	}

	return _address, nil
}

func (service *AddressService) Delete(ctx context.Context, ID primitive.ObjectID) error {
	return service.addressRepository.Delete(ctx, ID)
}
