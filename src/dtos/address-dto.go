package dtos

import "time"

type Address struct {
	ID        string    `json:"id"`
	Street    string    `json:"street,omitempty"`
	City      string    `json:"city,omitempty"`
	Province  string    `json:"province,omitempty"`
	Code      string    `json:"code,omitempty"`
	Type      string    `json:"type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Version   uint      `json:"version"`
}
