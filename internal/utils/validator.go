package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Optional: Add custom messages here. Format: "StructName.Field.Tag"
var customMessages = map[string]string{
	"CreatePostPayload.Title.required":   "Title is missing",
	"CreatePostPayload.Content.required": "Content is missing",
}

// ValidateStruct validates any struct and returns field-based error messages.
func ValidateStruct(data any) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	// Handle pointer to struct
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	structName := t.Name()

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fieldName := e.Field()
			tag := e.Tag()

			// Try to get custom message first
			key := fmt.Sprintf("%s.%s.%s", structName, fieldName, tag)
			if msg, ok := customMessages[key]; ok {
				errors[strings.ToLower(fieldName)] = msg
			} else {
				// Use generated message as fallback
				errors[strings.ToLower(fieldName)] = generateMessage(e)
			}
		}
	}

	return errors
}

// generateMessage creates user-friendly error messages based on tag and field type.
func generateMessage(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()
	fieldKind := e.Kind()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		switch fieldKind {
		case reflect.String:
			return fmt.Sprintf("%s must be at least %s characters long", field, param)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%s must be at least %s", field, param)
		case reflect.Slice, reflect.Array:
			return fmt.Sprintf("%s must contain at least %s items", field, param)
		default:
			return fmt.Sprintf("%s must be at least %s", field, param)
		}
	case "max":
		switch fieldKind {
		case reflect.String:
			return fmt.Sprintf("%s must be at most %s characters long", field, param)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%s must be at most %s", field, param)
		case reflect.Slice, reflect.Array:
			return fmt.Sprintf("%s must contain at most %s items", field, param)
		default:
			return fmt.Sprintf("%s must be at most %s", field, param)
		}
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "notempty":
		return fmt.Sprintf("%s cannot be empty", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
