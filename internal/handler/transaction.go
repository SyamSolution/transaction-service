package handler

import (
	"encoding/json"
	"github.com/SyamSolution/transaction-service/config"
	"github.com/SyamSolution/transaction-service/internal/model"
	"github.com/SyamSolution/transaction-service/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type transaction struct {
	transactionUsecase usecase.TransactionExecutor
	logger             config.Logger
}

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransactionByTransactionID(c *fiber.Ctx) error
	MidtransNotification(ctx *fiber.Ctx) error
	GetListTransaction(c *fiber.Ctx) error
}

func NewTransactionHandler(transactionUsecase usecase.TransactionExecutor, logger config.Logger) TransactionHandler {
	return &transaction{transactionUsecase: transactionUsecase, logger: logger}
}

func (h *transaction) CreateTransaction(c *fiber.Ctx) error {
	var request model.TransactionRequest
	if err := c.BodyParser(&request); err != nil {
		h.logger.Error("Error when parsing request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:3000/users/profile", nil)
	if err != nil {
		h.logger.Error("Error when creating new request", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	req.Header.Set("Authorization", c.Get("Authorization"))

	resp, err := client.Do(req)
	if err != nil {
		h.logger.Error("Error when sending request to user service", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("Error when reading response body", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	var respUser model.ResponseUser
	err = json.Unmarshal(body, &respUser)
	if err != nil {
		h.logger.Error("Error when unmarshalling response body", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	if respUser.User == (model.User{}) {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
			},
		})
	}

	snapResp, err := h.transactionUsecase.CreateTransaction(request, respUser.User)
	if err != nil {
		h.logger.Error("Error when creating transaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Data: struct {
			Token       string `json:"token"`
			RedirectURL string `json:"redirect_url"`
		}{
			Token:       snapResp.Token,
			RedirectURL: snapResp.RedirectURL,
		},
		Meta: model.Meta{
			Code:    fiber.StatusCreated,
			Message: "Transaction created successfully",
		},
	})
}

func (h *transaction) GetTransactionByTransactionID(c *fiber.Ctx) error {
	transactionID, err := c.ParamsInt("transaction_id")
	if err != nil {
		h.logger.Error("Error when parsing transaction ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}

	transaction, err := h.transactionUsecase.GetTransactionByTransactionID(transactionID)
	if err != nil {
		h.logger.Error("Error when getting transaction by transaction ID", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}
	transaction.Email = c.Locals("email").(string)

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Data: transaction,
		Meta: model.Meta{
			Code:    fiber.StatusOK,
			Message: "Transaction retrieved successfully",
		},
	})
}

func (h *transaction) GetListTransaction(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	var request model.TransactionListRequest
	if err := c.BodyParser(&request); err != nil {
		h.logger.Error("Error when parsing request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			},
		})
	}
	request.Email = email

	transactions, err := h.transactionUsecase.GetListTransaction(request)
	if err != nil {
		h.logger.Error("Error when getting list transaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Data: transactions,
		Meta: model.Meta{
			Code:    fiber.StatusOK,
			Message: "List transaction retrieved successfully",
		},
	})
}

func (h *transaction) MidtransNotification(ctx *fiber.Ctx) error {
	c := coreapi.Client{}
	c.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// 1. Initialize empty map
	var notificationPayload map[string]interface{}

	// 2. Parse JSON request body and use it to set json to payload
	if err := ctx.BodyParser(&notificationPayload); err != nil {
		// Return an error response if parsing fails
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// 3. Get order-id from payload
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		// Return an error response if order_id is not found in the payload
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order_id not found in request payload",
		})
	}

	// 4. Check transaction status using the orderId
	transactionStatusResp, err := c.CheckTransaction(orderId)
	if err != nil {
		// Return an error response if the transaction status check fails
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.GetMessage(),
		})
	}

	transaction, errors := h.transactionUsecase.GetTransactionByOrderID(orderId)
	if err != nil {
		h.logger.Error("Error when getting transaction by order ID", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.Error(),
		})
	}

	// 5. Set transaction status based on the response from the transaction status check
	if transactionStatusResp != nil {
		switch transactionStatusResp.TransactionStatus {
		case "capture":
			if transactionStatusResp.FraudStatus == "challenge" {
				// TODO: Update your database to set the transaction status to 'challenge'
			} else if transactionStatusResp.FraudStatus == "accept" {
				// TODO: Update your database to set the transaction status to 'success'
				err := h.transactionUsecase.UpdateTransactionStatus(orderId, "completed", transaction.Email)
				if err != nil {
					h.logger.Error("Error when updating transaction status", zap.Error(err))
					return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
			}
		case "settlement":
			// TODO: Update your database to set the transaction status to 'success'
			err := h.transactionUsecase.UpdateTransactionStatus(orderId, "completed", transaction.Email)
			if err != nil {
				h.logger.Error("Error when updating transaction status", zap.Error(err))
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		case "deny":
			// TODO: Handle 'deny' appropriately
		case "cancel", "expire":
			// TODO: Update your database to set the transaction status to 'failure'
			err := h.transactionUsecase.UpdateTransactionStatus(orderId, "cancelled", transaction.Email)
			if err != nil {
				h.logger.Error("Error when updating transaction status", zap.Error(err))
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		case "pending":
			// TODO: Update your database to set the transaction status to 'pending'
			err := h.transactionUsecase.UpdateTransactionStatus(orderId, "pending", transaction.Email)
			if err != nil {
				h.logger.Error("Error when updating transaction status", zap.Error(err))
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
	}

	// Return a success response
	return ctx.JSON(fiber.Map{
		"status": "ok",
	})
}
