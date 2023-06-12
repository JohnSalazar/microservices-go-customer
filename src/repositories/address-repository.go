package repositories

import (
	"context"
	"customer/src/models"
	"time"

	"github.com/JohnSalazar/microservices-go-common/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type addressRepository struct {
	database *mongo.Database
}

func NewAddressRepository(
	database *mongo.Database,
) *addressRepository {
	return &addressRepository{
		database: database,
	}
}

func (r *addressRepository) collectionName() string {
	return "addresses"
}

func (r *addressRepository) collection() *mongo.Collection {
	return r.database.Collection(r.collectionName())
}

func (r *addressRepository) find(ctx context.Context, filter interface{}) ([]*models.Address, error) {
	findOptions := options.FindOptions{}
	findOptions.SetSort(bson.M{"version": -1})

	newFilter := map[string]interface{}{
		"deleted": false,
	}
	mergeFilter := helpers.MergeFilters(newFilter, filter)

	cursor, err := r.collection().Find(ctx, mergeFilter, &findOptions)
	if err != nil {
		defer cursor.Close(ctx)
		return nil, err
	}

	var addresses []*models.Address

	for cursor.Next(ctx) {
		address := &models.Address{}

		cursor.Decode(address)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *addressRepository) findOne(ctx context.Context, filter interface{}) (*models.Address, error) {
	findOneOptions := options.FindOneOptions{}
	findOneOptions.SetSort(bson.M{"version": -1})

	newFilter := map[string]interface{}{
		"deleted": false,
	}
	mergeFilter := helpers.MergeFilters(newFilter, filter)

	address := &models.Address{}
	err := r.collection().FindOne(ctx, mergeFilter, &findOneOptions).Decode(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (r *addressRepository) findOneAndUpdate(ctx context.Context, filter interface{}, fields interface{}) *mongo.SingleResult {
	findOneAndUpdateOptions := options.FindOneAndUpdateOptions{}
	findOneAndUpdateOptions.SetReturnDocument(options.After)

	result := r.collection().FindOneAndUpdate(ctx, filter, bson.M{"$set": fields}, &findOneAndUpdateOptions)

	return result
}

func (r *addressRepository) GetByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*models.Address, error) {
	filter := bson.M{"customer_id": customerID}

	return r.find(ctx, filter)
}

func (r *addressRepository) FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Address, error) {
	filter := bson.M{"_id": ID}

	return r.findOne(ctx, filter)
}

func (r *addressRepository) FindAddressExists(ctx context.Context, address *models.Address) bool {
	filter := bson.M{"customer_id": address.CustomerID, "code": address.Code, "deleted": false}

	result := r.collection().FindOne(ctx, filter)

	return result.Err() == nil
}

func (r *addressRepository) Create(ctx context.Context, address *models.Address) error {
	fields := bson.M{
		"_id":         address.ID,
		"customer_id": address.CustomerID,
		"street":      address.Street,
		"city":        address.City,
		"province":    address.Province,
		"code":        address.Code,
		"type":        address.Type,
		"created_at":  time.Now().UTC(),
		"version":     0,
		"deleted":     false,
	}

	_, err := r.collection().InsertOne(ctx, fields)

	return err
}

func (r *addressRepository) Update(ctx context.Context, address *models.Address) (*models.Address, error) {
	address.Version++
	address.UpdatedAt = time.Now().UTC()
	fields := bson.M{
		"street":     address.Street,
		"city":       address.City,
		"province":   address.Province,
		"code":       address.Code,
		"type":       address.Type,
		"updated_at": address.UpdatedAt,
		"version":    address.Version,
	}

	filter := r.filterUpdate(address)

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return nil, result.Err()
	}

	modelAddress := &models.Address{}
	decodeErr := result.Decode(modelAddress)

	return modelAddress, decodeErr
}

func (r *addressRepository) Delete(ctx context.Context, ID primitive.ObjectID) error {
	filter := bson.M{"_id": ID}

	fields := bson.M{"deleted": true}

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *addressRepository) filterUpdate(address *models.Address) interface{} {
	filter := bson.M{
		"_id":     address.ID,
		"version": address.Version - 1,
	}

	return filter
}
