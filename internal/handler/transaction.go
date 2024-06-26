package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SyamSolution/transaction-service/config"
	"github.com/SyamSolution/transaction-service/helper"
	"github.com/SyamSolution/transaction-service/internal/model"
	"github.com/SyamSolution/transaction-service/internal/usecase"
	"github.com/SyamSolution/transaction-service/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"

	"go.uber.org/zap"
)

type transaction struct {
	transactionUsecase usecase.TransactionExecutor
	logger             config.Logger
	cacher             config.Cacher
}

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransactionByTransactionID(c *fiber.Ctx) error
	MidtransNotification(ctx *fiber.Ctx) error
	GetListTransaction(c *fiber.Ctx) error
	MidtransTransactionCancel(c *fiber.Ctx) error
}

func NewTransactionHandler(transactionUsecase usecase.TransactionExecutor, logger config.Logger, cacher config.Cacher) TransactionHandler {
	return &transaction{transactionUsecase: transactionUsecase, logger: logger, cacher: cacher}
}

func (h *transaction) CreateTransaction(c *fiber.Ctx) error {
	var request model.TransactionRequest
	if err := c.BodyParser(&request); err != nil {
		h.logger.Error("Error when parsing request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusBadRequest,
				Message: util.ERROR_NOT_FOUND_MSG,
			},
		})
	}

	//cek redis kalau ada
	var user model.User
	cacheDataUser, err := h.cacher.Get(c.Context(), "user-"+c.Locals("email").(string))
	if err != nil {
		h.logger.Error("Error when getting data from cache", zap.Error(err))
	}
	if cacheDataUser == "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", os.Getenv("USER_SERVICE_URL"), nil)
		if err != nil {
			h.logger.Error("Error when creating new request", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
				Meta: model.Meta{
					Code:    fiber.StatusInternalServerError,
					Message: util.ERROR_BASE_MSG,
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
					Message: util.ERROR_BASE_MSG,
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
					Message: util.ERROR_BASE_MSG,
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
					Message: util.ERROR_BASE_MSG,
				},
			})
		}

		if respUser.Data.User == (model.User{}) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ResponseWithoutData{
				Meta: model.Meta{
					Code:    fiber.StatusUnauthorized,
					Message: "Unauthorized",
				},
			})
		}
		user = respUser.Data.User

		userJson, _ := json.Marshal(respUser.Data.User)
		err = h.cacher.Set(c.Context(), "user-"+c.Locals("email").(string), userJson, time.Hour*24*30)
		if err != nil {
			h.logger.Error("Error when setting data to cache", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
				Meta: model.Meta{
					Code:    fiber.StatusInternalServerError,
					Message: util.ERROR_BASE_MSG,
				},
			})
		}
	} else {
		err = json.Unmarshal([]byte(cacheDataUser), &user)
		if err != nil {
			h.logger.Error("Error when unmarshalling response body", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
				Meta: model.Meta{
					Code:    fiber.StatusInternalServerError,
					Message: util.ERROR_BASE_MSG,
				},
			})
		}
	}

	snapResp, discount, total, err := h.transactionUsecase.CreateTransaction(request, user)
	if err != nil {
		if strings.Contains(err.Error(), "not eligible") {
			return c.Status(fiber.StatusBadRequest).JSON(model.ResponseWithoutData{
				Meta: model.Meta{
					Code:    fiber.StatusBadRequest,
					Message: "Not eligible to buy right now",
				},
			})
		} else {
			h.logger.Error("Error when creating transaction", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
				Meta: model.Meta{
					Code:    fiber.StatusInternalServerError,
					Message: util.ERROR_BASE_MSG,
				},
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Data: struct {
			Token            string  `json:"token"`
			RedirectURL      string  `json:"redirect_url"`
			Discount         float32 `json:"discount"`
			TotalTransaction float32 `json:"total_transaction"`
		}{
			Token:            snapResp.Token,
			RedirectURL:      snapResp.RedirectURL,
			Discount:         discount,
			TotalTransaction: total,
		},
		Meta: model.Meta{
			Code:    fiber.StatusCreated,
			Message: "Transaction created successfully",
		},
	})
}

func (h *transaction) GetTransactionByTransactionID(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	transactionID, err := c.ParamsInt("transaction_id")
	if err != nil {
		h.logger.Error("Error when parsing transaction ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusBadRequest,
				Message: util.ERROR_NOT_FOUND_MSG,
			},
		})
	}

	// TODO cek redis kalau ada data user
	transaction, err := h.transactionUsecase.GetTransactionByTransactionID(transactionID, email)
	if err != nil {
		h.logger.Error("Error when getting transaction by transaction ID", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: util.ERROR_BASE_MSG,
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
				Message: util.ERROR_NOT_FOUND_MSG,
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
				Message: util.ERROR_BASE_MSG,
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

	transactionOrder, errors := h.transactionUsecase.GetTransactionByOrderID(orderId)
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
				for _, dt := range transactionOrder.DetailTransactionResponse {
					message := model.MessageOrderTicket{
						TicketID: dt.TicketID,
						Order:    dt.Quantity,
					}

					jsonString, err := json.Marshal(message)
					if err != nil {
						fmt.Println("Error:", err)
						h.logger.Error("Error when proceesing message", zap.Error(err))
					
					}
					
					if err = helper.ProduceMessageSqs(os.Getenv("SQS_TICKET_SUCCESS_URL"), string(jsonString), "Update Ticket Success"); err != nil {
						h.logger.Error("Error when producing message", zap.Error(err))
					}
				}

				err := h.transactionUsecase.UpdateTransactionStatus(orderId, "completed", transactionOrder.Email)
				if err != nil {
					h.logger.Error("Error when updating transaction status", zap.Error(err))
					return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
			}
		case "settlement":
			// TODO: Update your database to set the transaction status to 'success'
			ticketEvent, err := helper.GetTicketEventByTicketID(transactionOrder.DetailTransactionResponse[0].TicketID)
			if err != nil {
				h.logger.Error("Error when getting ticket event by ticket ID", zap.Error(err))
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			emailPDF := model.EmailPDFMessage{
				Email:          transactionOrder.Email,
				OrderId:        orderId,
				EventName:      ticketEvent.EventName,
				Price:          transactionOrder.TotalAmount,
				NumberOfTicket: transactionOrder.TotalTicket,
				EventDate:      ticketEvent.Date.Format("2006-01-02"),
				EventTime:      ticketEvent.Date.Format("15:04:05"),
				Venue:          ticketEvent.CountryPlace,
				CustomerName:   transactionOrder.FullName,
				PurchaseDate:   transactionOrder.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			for _, dt := range transactionOrder.DetailTransactionResponse {
				emailPDF.DetailTickets = append(emailPDF.DetailTickets, model.DetailTicket{
					TicketType:  dt.TicketType,
					TotalTicket: dt.Quantity,
				})
			}
			jsonString, err := json.Marshal(emailPDF)
			if err != nil {
				fmt.Println("Error:", err)
				h.logger.Error("Error when proceesing message", zap.Error(err))
			
			}
			
			if err = helper.ProduceMessageSqs(os.Getenv("SQS_MAIL_URL"), string(jsonString), "Send Email PDF Success"); err != nil {
				h.logger.Error("Error when producing message", zap.Error(err))
			}

			for _, dt := range transactionOrder.DetailTransactionResponse {
				message := model.MessageOrderTicket{
					TicketID: dt.TicketID,
					Order:    dt.Quantity,
				}

				jsonString, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Error:", err)
					h.logger.Error("Error when proceesing message", zap.Error(err))
				
				}
				
				if err = helper.ProduceMessageSqs(os.Getenv("SQS_TICKET_SUCCESS_URL"), string(jsonString), "Update Ticket Success"); err != nil {
					h.logger.Error("Error when producing message", zap.Error(err))
				}
			}

			err = h.transactionUsecase.UpdateTransactionStatus(orderId, "completed", transactionOrder.Email)
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
			// kirim message SQS ke ticket-management-service balikin ticket
			for _, dt := range transactionOrder.DetailTransactionResponse {
				message := model.MessageOrderTicket{
					TicketID: dt.TicketID,
					Order:    dt.Quantity,
				}

				jsonString, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Error:", err)
					h.logger.Error("Error when proceesing message", zap.Error(err))
				
				}
				
				if err = helper.ProduceMessageSqs(os.Getenv("SQS_TICKET_FAILED_URL"), string(jsonString), "Update Ticket Failed"); err != nil {
					h.logger.Error("Error when producing message", zap.Error(err))
				}
			}

			err := h.transactionUsecase.UpdateTransactionStatus(orderId, "cancelled", transactionOrder.Email)
			if err != nil {
				h.logger.Error("Error when updating transaction status", zap.Error(err))
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		case "pending":
			// TODO: Update your database to set the transaction status to 'pending'
			err := h.transactionUsecase.UpdateTransactionStatus(orderId, "pending", transactionOrder.Email)
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

func (h *transaction) MidtransTransactionCancel(c *fiber.Ctx) error {
	orderId := c.Params("order_id")

	var core = coreapi.Client{}
	core.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	_, err := core.CancelTransaction(orderId)
	if err != nil {
		h.logger.Error("Error when cancelling transaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseWithoutData{
			Meta: model.Meta{
				Code:    fiber.StatusInternalServerError,
				Message: "Error when cancelling transaction",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.ResponseWithoutData{
		Meta: model.Meta{
			Code:    fiber.StatusOK,
			Message: "Transaction cancelled successfully",
		},
	})
}
