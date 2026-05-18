package courier

import (
	"context"
	"errors"
	"net/url"
	"time"

	"courier-technical-test/go/internal/response"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service   *Service
	validator *Validator
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service, validator: NewValidator()}
}

func (handler *Handler) Index(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	values, _ := url.ParseQuery(string(c.Context().URI().QueryString()))
	result, validationErrors := handler.service.List(ctx, values)
	if len(validationErrors) > 0 {
		return response.Failure(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}
	return response.List(c, "Couriers retrieved successfully", result.Couriers, response.Pagination{
		Page: result.Page, PerPage: result.PerPage, Total: result.Total, TotalPages: result.TotalPages,
	})
}

func (handler *Handler) Store(c *fiber.Ctx) error {
	var payload CourierPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Failure(c, fiber.StatusUnprocessableEntity, "Validation failed", map[string][]string{"body": []string{"Invalid JSON body."}})
	}
	if validationErrors := handler.validator.ValidatePayload(payload); len(validationErrors) > 0 {
		return response.Failure(c, fiber.StatusUnprocessableEntity, "Validation failed", validationErrors)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	courier, err := handler.service.Create(ctx, payload)
	if err != nil {
		return response.Failure(c, fiber.StatusUnprocessableEntity, "Validation failed", map[string][]string{"registered_at": []string{"The registered_at field must be a valid date or datetime."}})
	}
	return response.Success(c, fiber.StatusCreated, "Courier created successfully", courier)
}

func (handler *Handler) Show(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	courier, err := handler.service.Find(ctx, c.Params("id"))
	if errors.Is(err, ErrInvalidID) || errors.Is(err, ErrNotFound) {
		return response.Failure(c, fiber.StatusNotFound, "Courier not found", nil)
	}
	if err != nil {
		return response.Failure(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}
	return response.Success(c, fiber.StatusOK, "Courier retrieved successfully", courier)
}

func (handler *Handler) Update(c *fiber.Ctx) error {
	var payload CourierPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Failure(c, fiber.StatusUnprocessableEntity, "Validation failed", map[string][]string{"body": []string{"Invalid JSON body."}})
	}
	if validationErrors := handler.validator.ValidatePayload(payload); len(validationErrors) > 0 {
		return response.Failure(c, fiber.StatusUnprocessableEntity, "Validation failed", validationErrors)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	courier, err := handler.service.Update(ctx, c.Params("id"), payload)
	if errors.Is(err, ErrInvalidID) || errors.Is(err, ErrNotFound) {
		return response.Failure(c, fiber.StatusNotFound, "Courier not found", nil)
	}
	if err != nil {
		return response.Failure(c, fiber.StatusUnprocessableEntity, "Validation failed", map[string][]string{"registered_at": []string{"The registered_at field must be a valid date or datetime."}})
	}
	return response.Success(c, fiber.StatusOK, "Courier updated successfully", courier)
}

func (handler *Handler) Destroy(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := handler.service.Delete(ctx, c.Params("id"))
	if errors.Is(err, ErrInvalidID) || errors.Is(err, ErrNotFound) {
		return response.Failure(c, fiber.StatusNotFound, "Courier not found", nil)
	}
	if err != nil {
		return response.Failure(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}
	return response.Success(c, fiber.StatusOK, "Courier deleted successfully", nil)
}
