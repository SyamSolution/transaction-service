package handler

import (
	"encoding/json"
	"fmt"
	"github.com/SyamSolution/transaction-service/config"
	"github.com/SyamSolution/transaction-service/internal/model"
	"github.com/SyamSolution/transaction-service/internal/usecase"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"os"

	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"

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
}

func NewTransactionHandler(transactionUsecase usecase.TransactionExecutor, logger config.Logger) TransactionHandler {
	return &transaction{transactionUsecase: transactionUsecase, logger: logger}
}

func (h *transaction) CreateTransaction(c *fiber.Ctx) error {
	var request model.TransactionRequest
	if err := c.BodyParser(&request); err != nil {
		h.logger.Error("Error when parsing request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:3000/users/profile", nil)
	if err != nil {
		h.logger.Error("Error when creating new request", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when creating new request",
		})
	}

	req.Header.Set("Authorization", c.Get("Authorization"))

	resp, err := client.Do(req)
	if err != nil {
		h.logger.Error("Error when sending request to user service", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when sending request to user service",
		})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("Error when reading response body", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when reading response body",
		})
	}

	var respUser model.ResponseUser
	err = json.Unmarshal(body, &respUser)
	if err != nil {
		h.logger.Error("Error when unmarshalling response body", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when unmarshalling response body",
		})
	}

	if respUser.User == (model.User{}) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	snapResp, err := h.transactionUsecase.CreateTransaction(request, respUser.User)
	if err != nil {
		h.logger.Error("Error when creating transaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	//ngirim ke notification-service untuk mengirimkan email data detail transaksi

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Transaction created successfully",
		"token":        snapResp.Token,
		"redirect_url": snapResp.RedirectURL,
	})
}

func (h *transaction) GetTransactionByTransactionID(c *fiber.Ctx) error {
	transactionID, err := c.ParamsInt("transaction_id")
	if err != nil {
		h.logger.Error("Error when parsing transaction ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid transaction ID",
		})
	}

	transaction, err := h.transactionUsecase.GetTransactionByTransactionID(transactionID)
	if err != nil {
		h.logger.Error("Error when getting transaction by transaction ID", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	transaction.Email = c.Locals("email").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transaction": transaction,
	})
}

func (h *transaction) GetListTransaction(c *fiber.Ctx) error {
	//email := c.Locals("email").(string)
	//status := c.Query("status")

	//request := model.TransactionListRequest{
	//	Email:  email,
	//	Status: status,
	//}
	//
	//transactions, err := h.transactionUsecase.GetListTransaction(request)
	//if err != nil {
	//	h.logger.Error("Error when getting list transaction", zap.Error(err))
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//		"message": err.Error(),
	//	})
	//}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		//"transactions": transactions,
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

	fmt.Println(orderId)

	// 4. Check transaction status using the orderId
	transactionStatusResp, err := c.CheckTransaction(orderId)
	if err != nil {
		// Return an error response if the transaction status check fails
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.GetMessage(),
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
				fmt.Printf("transaction success")
			}
		case "settlement":
			// TODO: Update your database to set the transaction status to 'success'
			err := h.transactionUsecase.UpdateTransactionStatus(orderId, "completed")
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
		case "pending":
			// TODO: Update your database to set the transaction status to 'pending'
		}
	}

	// Return a success response
	return ctx.JSON(fiber.Map{
		"status": "ok",
	})
}
