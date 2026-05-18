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
	ctx, cancel := requestContext()
	defer cancel()

	values, _ := url.ParseQuery(string(c.Context().URI().QueryString()))
	params, validationErrors := handler.validator.ParseQuery(values)
	if len(validationErrors) > 0 {
		return response.Failure(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}
	result, err := handler.service.List(ctx, params)
	if err != nil {
		return response.Failure(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}
	return response.List(c, "Couriers retrieved successfully", result.Couriers, response.Pagination{
		Page: result.Page, PerPage: result.PerPage, Total: result.Total, TotalPages: result.TotalPages,
	})
}

func (handler *Handler) Store(c *fiber.Ctx) error {
	var payload CourierPayload
	if err := c.BodyParser(&payload); err != nil {
		return validationFailure(c, map[string][]string{"body": {"Invalid JSON body."}})
	}
	if validationErrors := handler.validator.ValidatePayload(payload); len(validationErrors) > 0 {
		return validationFailure(c, validationErrors)
	}

	ctx, cancel := requestContext()
	defer cancel()
	courier, err := handler.service.Create(ctx, payload)
	if errors.Is(err, ErrInvalidRegisteredAt) {
		return invalidRegisteredAt(c)
	}
	if err != nil {
		return response.Failure(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}
	return response.Success(c, fiber.StatusCreated, "Courier created successfully", courier)
}

func (handler *Handler) Show(c *fiber.Ctx) error {
	ctx, cancel := requestContext()
	defer cancel()
	courier, err := handler.service.Find(ctx, c.Params("id"))
	if err != nil {
		return lookupFailure(c, err)
	}
	return response.Success(c, fiber.StatusOK, "Courier retrieved successfully", courier)
}

func (handler *Handler) Update(c *fiber.Ctx) error {
	var payload CourierPayload
	if err := c.BodyParser(&payload); err != nil {
		return validationFailure(c, map[string][]string{"body": {"Invalid JSON body."}})
	}
	if validationErrors := handler.validator.ValidatePayload(payload); len(validationErrors) > 0 {
		return validationFailure(c, validationErrors)
	}

	ctx, cancel := requestContext()
	defer cancel()
	courier, err := handler.service.Update(ctx, c.Params("id"), payload)
	if errors.Is(err, ErrInvalidID) || errors.Is(err, ErrNotFound) {
		return response.Failure(c, fiber.StatusNotFound, "Courier not found", nil)
	}
	if errors.Is(err, ErrInvalidRegisteredAt) {
		return invalidRegisteredAt(c)
	}
	if err != nil {
		return response.Failure(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}
	return response.Success(c, fiber.StatusOK, "Courier updated successfully", courier)
}

func (handler *Handler) Destroy(c *fiber.Ctx) error {
	ctx, cancel := requestContext()
	defer cancel()
	err := handler.service.Delete(ctx, c.Params("id"))
	if err != nil {
		return lookupFailure(c, err)
	}
	return response.Success(c, fiber.StatusOK, "Courier deleted successfully", nil)
}

func requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func validationFailure(c *fiber.Ctx, errorsByField map[string][]string) error {
	return response.Failure(c, fiber.StatusUnprocessableEntity, "Validation failed", errorsByField)
}

func lookupFailure(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrInvalidID) || errors.Is(err, ErrNotFound) {
		return response.Failure(c, fiber.StatusNotFound, "Courier not found", nil)
	}
	return response.Failure(c, fiber.StatusInternalServerError, "Internal server error", nil)
}

func invalidRegisteredAt(c *fiber.Ctx) error {
	return validationFailure(c, map[string][]string{"registered_at": {"The registered_at field must be a valid date or datetime."}})
}
