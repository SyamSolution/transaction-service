package util

import (
	"net/http"

	"github.com/SyamSolution/transaction-service/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(code int, data interface{}) model.Response {
	return model.Response{
		Data: data,
		Meta: model.Meta{
			Code:    code,
			Message: SUCCESS_RESPONSE_MSG,
		},
	}
}

func SuccessResponseWithoutData(code int, message string) model.ResponseWithoutData {
	return model.ResponseWithoutData{
		Meta: model.Meta{
			Code:    code,
			Message: message,
		},
	}
}

func BusinessErrorResponseWithoutData(err *model.BusinessError) model.ResponseWithoutData {
	code := DEFAULT_BUSINESS_ERROR_CODE
	message := DEFAULT_BUSINESS_ERROR_MESSAGE
	if err != nil {
		code = err.Code
		message = err.Message
	}
	return model.ResponseWithoutData{
		Meta: model.Meta{
			Code:    code,
			Message: message,
		},
	}
}

func BusinessErrorResponse(err *model.BusinessError) model.Response {
	code := DEFAULT_BUSINESS_ERROR_CODE
	message := DEFAULT_BUSINESS_ERROR_MESSAGE
	if err != nil {
		code = err.Code
		message = err.Message
	}
	return model.Response{
		Data: nil,
		Meta: model.Meta{
			Code:    code,
			Message: message,
		},
	}
}

func ValidatorErrorResponse(err []*model.ErrorFieldResponse) model.Response {
	return model.Response{
		Data: nil,
		Meta: model.Meta{
			Code:    ERROR_INVALID_PARAM_CODE,
			Message: ERROR_INVALID_PARAM_MSG,
		},
		Errors: err,
	}
}

func ValidateStruct(err error) []*model.ErrorFieldResponse {
	var errors []*model.ErrorFieldResponse
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element model.ErrorFieldResponse
			element.Field = err.Field()
			element.ErrMessage = err.Tag()
			errors = append(errors, &element)
		}
	}
	return errors
}

func AbortUnauthorized(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusUnauthorized).JSON(model.Response{
		Data: nil,
		Meta: model.Meta{
			Code:    ERROR_UNAUTHORIZE_CODE,
			Message: ERROR_UNAUTHORIZE_MSG,
		},
	})
}

func InternalBusinessError() *model.BusinessError {
	return &model.BusinessError{
		Code:    DEFAULT_BUSINESS_ERROR_CODE,
		Message: DEFAULT_BUSINESS_ERROR_MESSAGE,
	}
}
