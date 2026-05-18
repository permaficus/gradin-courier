package courier

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
}

func (v *Validator) ValidatePayload(payload CourierPayload) map[string][]string {
	if err := v.validate.Struct(payload); err != nil {
		errorsByField := map[string][]string{}
		for _, fieldError := range err.(validator.ValidationErrors) {
			field := jsonFieldName(fieldError.Field())
			errorsByField[field] = []string{fmt.Sprintf("The %s field is invalid.", field)}
		}
		return errorsByField
	}
	return nil
}

func (v *Validator) ParseQuery(raw url.Values) (QueryParams, map[string][]string) {
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

func jsonFieldName(field string) string {
	names := map[string]string{
		"Name": "name", "Email": "email", "Phone": "phone", "Level": "level",
		"VehicleType": "vehicle_type", "LicensePlate": "license_plate",
		"Status": "status", "RegisteredAt": "registered_at",
	}
	if name, ok := names[field]; ok {
		return name
	}
	return field
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
