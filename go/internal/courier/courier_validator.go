package courier

import (
	"fmt"

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
