package routes

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/handlers"
	"merchants.sidooh/pkg/services/transaction"
)

func TransactionRouter(app fiber.Router, service transaction.Service) {
	app.Get("/transactions", handlers.GetTransactions(service))
	app.Get("/transactions/:id", handlers.GetTransaction(service))

	app.Get("/merchants/:merchantId/transactions", handlers.GetTransactionsByMerchant(service))
	app.Post("/merchants/:merchantId/buy-float", handlers.BuyFloat(service))
	app.Post("/merchants/:merchantId/earnings/withdraw", handlers.WithdrawEarnings(service))

	app.Post("/payments/ipn", handlers.BuyFloat(service))
}
