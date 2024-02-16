package middleware

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/utils"
)

var Validator *validator.Validate

func validate(i interface{}) error {
	if err := Validator.Struct(i); err != nil {
		fmt.Println(err)
		var validationErrors []utils.ValidationError

		for _, err := range err.(validator.ValidationErrors) {

			msg := fmt.Sprintf("%v is invalid/missing", err.Field())

			tag := err.Tag()
			if err.Param() != "" {
				tag += " " + err.Param()
			}

			validationErrors = append(validationErrors, utils.ValidationError{
				Value:   err.Value(),
				Field:   err.Field(),
				Message: msg,
				Param:   tag,
			})
		}

		// Optionally, you could return the error to give each route more control over the status code
		return utils.ValidationErrorResponse(validationErrors)
	}
	return nil
}

func BindAndValidateRequest(context *fiber.Ctx, request interface{}) error {
	if err := context.BodyParser(request); err != nil {
		return err
	}

	if err := validate(request); err != nil {
		return err
	}

	return nil
}
