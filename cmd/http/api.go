package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SyamSolution/transaction-service/config"
	"github.com/SyamSolution/transaction-service/config/middleware"
	"github.com/SyamSolution/transaction-service/internal/handler"
	"github.com/SyamSolution/transaction-service/internal/repository"
	"github.com/SyamSolution/transaction-service/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	baseDep := config.NewBaseDep()
	loadEnv(baseDep.Logger)
	db, err := config.NewDbPool(baseDep.Logger)
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	cacher := config.NewCacher(baseDep.Logger)

	dbCollector := middleware.NewStatsCollector("assesment", db)
	prometheus.MustRegister(dbCollector)
	fiberProm := middleware.NewWithRegistry(prometheus.DefaultRegisterer, "transaction-service", "", "", map[string]string{})

	//=== repository lists start ===//
	transactionRepo := repository.NewTransactionRepository(db, baseDep.Logger)
	//=== repository lists end ===//

	//=== usecase lists start ===//
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, baseDep.Logger)
	//=== usecase lists end ===//

	//=== handler lists start ===//
	transactionHandler := handler.NewTransactionHandler(transactionUsecase, baseDep.Logger, cacher)
	//=== handler lists end ===//

	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeInterval: time.Millisecond,
		TimeFormat:   "02-01-2006 15:04:05",
		TimeZone:     "Indonesia/Jakarta",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//=== metrics route
	fiberProm.RegisterAt(app, "/metrics")
	app.Use(fiberProm.Middleware)

	//=== transaction routes ===//
	app.Post("/midtrans-notification", transactionHandler.MidtransNotification)
	app.Group("/", middleware.Auth())
	app.Post("/transactions", transactionHandler.CreateTransaction)
	app.Get("/transactions/:transaction_id", transactionHandler.GetTransactionByTransactionID)
	app.Get("/transactions-list", transactionHandler.GetListTransaction)
	app.Post("/midtrans/transaction-cancel/:order_id", transactionHandler.MidtransTransactionCancel)

	//=== listen port ===//
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
		log.Fatal(err)
	}
}

func loadEnv(logger config.Logger) {
	_, err := os.Stat(".env")
	if err == nil {
		err = godotenv.Load()
		if err != nil {
			logger.Error("no .env files provided")
		}
	}
}
