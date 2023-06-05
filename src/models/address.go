package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customerId,omitempty"`
	Street     string             `bson:"street" json:"street,omitempty"`
	City       string             `bson:"city" json:"city,omitempty"`
	Province   string             `bson:"province" json:"province,omitempty"`
	Code       string             `bson:"code" json:"code,omitempty"`
	Type       string             `bson:"type" json:"type,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	Version    uint               `bson:"version" json:"version,omitempty"`
	Deleted    bool               `bson:"deleted" json:"deleted,omitempty"`
}
