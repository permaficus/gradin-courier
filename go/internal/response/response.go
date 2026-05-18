package response

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Meta struct {
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

type SuccessBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    Meta   `json:"meta"`
}

type ListBody struct {
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
	Meta       Meta       `json:"meta"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type ErrorBody struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
	Meta    Meta                `json:"meta"`
}

func RequestID(c *fiber.Ctx) string {
	requestID := c.Get("x-request-id")
	if requestID == "" {
		requestID = newRequestID()
	}
	c.Set("x-request-id", requestID)
	return requestID
}

func newRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
}

func meta(c *fiber.Ctx) Meta {
	return Meta{RequestID: RequestID(c), Timestamp: time.Now().UTC()}
}

func Success(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(SuccessBody{Success: true, Message: message, Data: data, Meta: meta(c)})
}

func List(c *fiber.Ctx, message string, data any, pagination Pagination) error {
	return c.Status(fiber.StatusOK).JSON(ListBody{Success: true, Message: message, Data: data, Pagination: pagination, Meta: meta(c)})
}

func Failure(c *fiber.Ctx, status int, message string, errors map[string][]string) error {
	return c.Status(status).JSON(ErrorBody{Success: false, Message: message, Errors: errors, Meta: meta(c)})
}
