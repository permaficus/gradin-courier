package courier

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const CollectionName = "couriers"

type Courier struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string        `bson:"name" json:"name"`
	Email        *string       `bson:"email,omitempty" json:"email,omitempty"`
	Phone        *string       `bson:"phone,omitempty" json:"phone,omitempty"`
	Level        int           `bson:"level" json:"level"`
	VehicleType  *string       `bson:"vehicle_type,omitempty" json:"vehicle_type,omitempty"`
	LicensePlate *string       `bson:"license_plate,omitempty" json:"license_plate,omitempty"`
	Status       string        `bson:"status" json:"status"`
	RegisteredAt time.Time     `bson:"registered_at" json:"registered_at"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time    `bson:"deleted_at" json:"deleted_at,omitempty"`
}

type CourierPayload struct {
	Name         string  `json:"name" validate:"required,min=2,max=150"`
	Email        *string `json:"email" validate:"omitempty,email,max=150"`
	Phone        *string `json:"phone" validate:"omitempty,max=30"`
	Level        int     `json:"level" validate:"required,min=1,max=5"`
	VehicleType  *string `json:"vehicle_type" validate:"omitempty,max=50"`
	LicensePlate *string `json:"license_plate" validate:"omitempty,max=30"`
	Status       *string `json:"status" validate:"omitempty,oneof=active inactive suspended"`
	RegisteredAt *string `json:"registered_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type QueryParams struct {
	Page    int
	PerPage int
	Sort    string
	Search  string
	Levels  []int
}

type ListResult struct {
	Couriers   []Courier
	Page       int
	PerPage    int
	Total      int64
	TotalPages int64
}
