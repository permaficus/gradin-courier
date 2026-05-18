package courier

import (
	"context"
	"errors"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrNotFound = errors.New("courier not found")
var ErrInvalidID = errors.New("invalid courier id")

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (service *Service) Create(ctx context.Context, payload CourierPayload) (Courier, error) {
	now := time.Now().UTC()
	registeredAt, err := parseRegisteredAt(payload.RegisteredAt, now)
	if err != nil {
		return Courier{}, err
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

func (service *Service) List(ctx context.Context, raw url.Values) (ListResult, map[string][]string) {
	params, validationErrors := ParseQuery(raw)
	if len(validationErrors) > 0 {
		return ListResult{}, validationErrors
	}
	result, err := service.repository.List(ctx, params)
	if err != nil {
		return ListResult{}, map[string][]string{"database": []string{err.Error()}}
	}
	return result, nil
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
		return Courier{}, err
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

func ParseQuery(raw url.Values) (QueryParams, map[string][]string) {
	errorsByField := map[string][]string{}
	page := intParam(raw.Get("page"), 1, 1, math.MaxInt, "page", errorsByField)
	perPage := intParam(raw.Get("per_page"), 10, 1, 100, "per_page", errorsByField)
	sort := raw.Get("sort")
	allowedSort := map[string]bool{"": true, "name": true, "-name": true, "registered_at": true, "-registered_at": true, "created_at": true, "-created_at": true}
	if !allowedSort[sort] {
		errorsByField["sort"] = []string{"The sort field is invalid."}
	}
	levels := []int{}
	if levelRaw := strings.TrimSpace(raw.Get("level")); levelRaw != "" {
		for _, part := range strings.Split(levelRaw, ",") {
			level, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || level < 1 || level > 5 {
				errorsByField["level"] = []string{"The level query must contain only levels 1 to 5."}
				break
			}
			levels = append(levels, level)
		}
	}
	return QueryParams{Page: page, PerPage: perPage, Sort: sort, Search: raw.Get("search"), Levels: levels}, errorsByField
}

func intParam(raw string, fallback int, min int, max int, field string, errorsByField map[string][]string) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		errorsByField[field] = []string{"The " + field + " query is invalid."}
		return fallback
	}
	return value
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
