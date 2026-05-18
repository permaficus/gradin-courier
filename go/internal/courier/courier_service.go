package courier

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrNotFound = errors.New("courier not found")
var ErrInvalidID = errors.New("invalid courier id")
var ErrInvalidRegisteredAt = errors.New("invalid registered_at")

type Service struct {
	repository Store
}

type Store interface {
	Create(ctx context.Context, courier Courier) (Courier, error)
	List(ctx context.Context, params QueryParams) (ListResult, error)
	FindByID(ctx context.Context, id bson.ObjectID) (Courier, error)
	Update(ctx context.Context, id bson.ObjectID, courier Courier) (Courier, error)
	SoftDelete(ctx context.Context, id bson.ObjectID, deletedAt time.Time) error
}

func NewService(repository Store) *Service {
	return &Service{repository: repository}
}

func (service *Service) Create(ctx context.Context, payload CourierPayload) (Courier, error) {
	now := time.Now().UTC()
	registeredAt, err := parseRegisteredAt(payload.RegisteredAt, now)
	if err != nil {
		return Courier{}, ErrInvalidRegisteredAt
	}
	status := "active"
	if payload.Status != nil && *payload.Status != "" {
		status = *payload.Status
	}
	return service.repository.Create(ctx, Courier{
		Name: payload.Name, Email: payload.Email, Phone: payload.Phone, Level: payload.Level,
		VehicleType: payload.VehicleType, LicensePlate: payload.LicensePlate, Status: status,
		RegisteredAt: registeredAt, CreatedAt: now, UpdatedAt: now, DeletedAt: nil,
	})
}

func (service *Service) List(ctx context.Context, params QueryParams) (ListResult, error) {
	return service.repository.List(ctx, params)
}

func (service *Service) Find(ctx context.Context, id string) (Courier, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Courier{}, ErrInvalidID
	}
	courier, err := service.repository.FindByID(ctx, objectID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Courier{}, ErrNotFound
	}
	return courier, err
}

func (service *Service) Update(ctx context.Context, id string, payload CourierPayload) (Courier, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Courier{}, ErrInvalidID
	}
	current, err := service.repository.FindByID(ctx, objectID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Courier{}, ErrNotFound
	}
	if err != nil {
		return Courier{}, err
	}
	registeredAt, err := parseRegisteredAt(payload.RegisteredAt, current.RegisteredAt)
	if err != nil {
		return Courier{}, ErrInvalidRegisteredAt
	}
	status := "active"
	if payload.Status != nil && *payload.Status != "" {
		status = *payload.Status
	}
	return service.repository.Update(ctx, objectID, Courier{
		Name: payload.Name, Email: payload.Email, Phone: payload.Phone, Level: payload.Level,
		VehicleType: payload.VehicleType, LicensePlate: payload.LicensePlate, Status: status,
		RegisteredAt: registeredAt, UpdatedAt: time.Now().UTC(),
	})
}

func (service *Service) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	if err := service.repository.SoftDelete(ctx, objectID, time.Now().UTC()); errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func parseRegisteredAt(raw *string, fallback time.Time) (time.Time, error) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return fallback, nil
	}
	if parsed, err := time.Parse(time.RFC3339, *raw); err == nil {
		return parsed.UTC(), nil
	}
	parsed, err := time.Parse("2006-01-02", *raw)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}
