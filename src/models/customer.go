package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email     string             `bson:"email" json:"email,omitempty"`
	FirstName string             `bson:"firstname" json:"firstname,omitempty"`
	LastName  string             `bson:"lastname" json:"lastname,omitempty"`
	Avatar    string             `bson:"avatar" json:"avatar,omitempty"`
	Phone     string             `bson:"phone" json:"phone,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	Version   uint               `bson:"version" json:"version,omitempty"`
	Deleted   bool               `bson:"deleted" json:"deleted,omitempty"`
}
