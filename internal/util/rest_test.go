package util

import (
	"testing"

	"github.com/SyamSolution/transaction-service/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestSuccessResponse(t *testing.T) {
	t.Run("should return standard model success response", func(t *testing.T) {
		data := "this is data"
		code := 2002

		actResp := SuccessResponse(code, data)
		assert.NotEmpty(t, actResp)
		assert.Equal(t, actResp.Data, data)
		assert.Equal(t, actResp.Meta.Code, code)
		assert.Equal(t, actResp.Meta.Message, "success")
	})
}

func TestBusinessErrorResponse(t *testing.T) {
	t.Run("should return model response contains business error message", func(t *testing.T) {
		businessError := model.BusinessError{
			Code:    4113,
			Message: "input tidak sesuai",
		}

		actResp := BusinessErrorResponse(&businessError)
		assert.NotEmpty(t, actResp)
		assert.Nil(t, actResp.Data)
		assert.Equal(t, actResp.Meta.Code, businessError.Code)
		assert.Equal(t, actResp.Meta.Message, businessError.Message)
	})

	t.Run("should return model response contains business error message nil", func(t *testing.T) {
		actResp := BusinessErrorResponse(nil)
		assert.NotEmpty(t, actResp)
		assert.Nil(t, actResp.Data)
		assert.Equal(t, actResp.Meta.Code, 4001)
		assert.Equal(t, actResp.Meta.Message, "something is wrong, report to support team")
	})
}

func TestValidatorErrorResponse(t *testing.T) {
	t.Run("should return standard validator error response", func(t *testing.T) {
		actValidatorRes := []*model.ErrorFieldResponse{
			{
				Field:      "id",
				ErrMessage: "is required",
			},
		}

		actResp := ValidatorErrorResponse(actValidatorRes)
		assert.NotEmpty(t, actResp)
		assert.Nil(t, actResp.Data)
		assert.Equal(t, actResp.Meta.Code, 4101)
		assert.Equal(t, actResp.Meta.Message, "failed")
		assert.Equal(t, actResp.Errors, actValidatorRes)
	})

	t.Run("should return standard validator error response with nil error", func(t *testing.T) {
		actResp := ValidatorErrorResponse(nil)
		assert.NotEmpty(t, actResp)
		assert.Nil(t, actResp.Data)
		assert.Equal(t, actResp.Meta.Code, 4101)
		assert.Equal(t, actResp.Meta.Message, "failed")
		assert.Nil(t, actResp.Errors)
	})
}

func TestAbortUnauthorized(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	t.Run("should return default unauthorized request", func(t *testing.T) {
		err := AbortUnauthorized(c)
		assert.Nil(t, err)
	})
}

func TestSuccessResponseWithoutData(t *testing.T) {
	t.Run("should return code and default success message", func(t *testing.T) {
		paramCode := 2001
		paramMessage := "success"
		res := SuccessResponseWithoutData(paramCode, paramMessage)

		assert.Equal(t, paramCode, res.Meta.Code)
		assert.Equal(t, paramMessage, res.Meta.Message)
	})
}

func TestInternalBusinessError(t *testing.T) {
	t.Run("should return code and default success message", func(t *testing.T) {
		res := InternalBusinessError()

		assert.Equal(t, 4001, res.Code)
		assert.Equal(t, "something is wrong, report to support team", res.Message)
	})
}

func TestBusinessErrorResponseWithoutData(t *testing.T) {
	t.Run("should return model response contains business error message without data", func(t *testing.T) {
		businessError := model.BusinessError{
			Code:    4113,
			Message: "input tidak sesuai",
		}

		actResp := BusinessErrorResponseWithoutData(&businessError)
		assert.NotEmpty(t, actResp)
		assert.Equal(t, actResp.Meta.Code, businessError.Code)
		assert.Equal(t, actResp.Meta.Message, businessError.Message)
	})

	t.Run("should return model response contains business error message nil", func(t *testing.T) {
		actResp := BusinessErrorResponseWithoutData(nil)
		assert.NotEmpty(t, actResp)
		assert.Equal(t, actResp.Meta.Code, 4001)
		assert.Equal(t, actResp.Meta.Message, "something is wrong, report to support team")
	})
}
