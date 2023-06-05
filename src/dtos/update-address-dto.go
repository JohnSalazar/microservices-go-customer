package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateAddress struct {
	ID       primitive.ObjectID `json:"id"`
	Street   string             `json:"street,omitempty"`
	City     string             `json:"city,omitempty"`
	Province string             `json:"province,omitempty"`
	Code     string             `json:"code,omitempty"`
	Type     string             `json:"type,omitempty"`
	Version  uint               `json:"version"`
}
