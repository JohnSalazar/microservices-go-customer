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

type customerRepository struct {
	database *mongo.Database
}

func NewCustomerRepository(
	database *mongo.Database,
) *customerRepository {
	return &customerRepository{
		database: database,
	}
}

func (r *customerRepository) collectionName() string {
	return "customers"
}

func (r *customerRepository) collection() *mongo.Collection {
	return r.database.Collection(r.collectionName())
}

func (r *customerRepository) findOne(ctx context.Context, filter interface{}) (*models.Customer, error) {
	findOneOptions := options.FindOneOptions{}
	findOneOptions.SetSort(bson.M{"version": -1})

	newFilter := map[string]interface{}{
		"deleted": false,
	}
	mergeFilter := helpers.MergeFilters(newFilter, filter)

	customer := &models.Customer{}
	err := r.collection().FindOne(ctx, mergeFilter, &findOneOptions).Decode(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *customerRepository) findOneAndUpdate(ctx context.Context, filter interface{}, fields interface{}) *mongo.SingleResult {
	findOneAndUpdateOptions := options.FindOneAndUpdateOptions{}
	findOneAndUpdateOptions.SetReturnDocument(options.After)

	result := r.collection().FindOneAndUpdate(ctx, filter, bson.M{"$set": fields}, &findOneAndUpdateOptions)

	return result
}

func (r *customerRepository) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	filter := bson.M{"email": email}

	return r.findOne(ctx, filter)
}

func (r *customerRepository) FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Customer, error) {
	filter := bson.M{"_id": ID}

	return r.findOne(ctx, filter)
}

func (r *customerRepository) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer.CreatedAt = time.Now().UTC()

	fields := bson.M{
		"_id":        customer.ID,
		"email":      customer.Email,
		"firstname":  customer.FirstName,
		"lastname":   customer.LastName,
		"avatar":     customer.Avatar,
		"phone":      customer.Phone,
		"created_at": customer.CreatedAt,
		"version":    0,
		"deleted":    false,
	}

	_, err := r.collection().InsertOne(ctx, fields)

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer.Version++
	customer.UpdatedAt = time.Now().UTC()
	fields := bson.M{
		"firstname":  customer.FirstName,
		"lastname":   customer.LastName,
		"avatar":     customer.Avatar,
		"phone":      customer.Phone,
		"updated_at": customer.UpdatedAt,
		"version":    customer.Version,
	}

	filter := r.filterUpdate(customer)

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return nil, result.Err()
	}

	modelCustomer := &models.Customer{}
	decodeErr := result.Decode(modelCustomer)

	return modelCustomer, decodeErr
}

func (r *customerRepository) Delete(ctx context.Context, ID primitive.ObjectID) error {
	filter := bson.M{"_id": ID}

	fields := bson.M{"deleted": true}

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *customerRepository) filterUpdate(customer *models.Customer) interface{} {
	filter := bson.M{
		"_id":     customer.ID,
		"version": customer.Version - 1,
	}

	return filter
}
