package api

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"merchants.sidooh/api/middleware"
	"merchants.sidooh/api/middleware/jwt"
	"merchants.sidooh/api/routes"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/earning"
	"merchants.sidooh/pkg/services/earning_account"
	"merchants.sidooh/pkg/services/earning_account_transaction"
	"merchants.sidooh/pkg/services/ipn"
	"merchants.sidooh/pkg/services/jobs"
	"merchants.sidooh/pkg/services/location"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/pkg/services/mpesa_store"
	"merchants.sidooh/pkg/services/payment"
	"merchants.sidooh/pkg/services/transaction"
	"merchants.sidooh/utils"
	"time"
)

func setMiddleware(app *fiber.App) {
	//app.Use(func(c *fiber.Ctx) error {
	//	// TODO: URGENT: Check out these headers and review them
	//	// Set some security headers:
	//	c.Set("X-XSS-Protection", "1; mode=block")
	//	c.Set("X-Content-Type-Options", "nosniff")
	//	c.Set("X-Download-Options", "noopen")
	//	c.Set("Strict-Transport-Security", "max-age=5184000")
	//	c.Set("X-Frame-Options", "SAMEORIGIN")
	//	c.Set("X-DNS-Prefetch-Control", "off")
	//
	//	// Go to next middleware:
	//	return c.Next()
	//})

	app.Use(helmet.New())
	app.Use(cors.New())

	app.Use(limiter.New(limiter.Config{Max: viper.GetInt("RATE_LIMIT")}))
	app.Use(recover.New())
	app.Use(fiberLogger.New(fiberLogger.Config{Output: utils.GetLogFile("stats.log")}))

	app.Use(favicon.New(favicon.Config{
		Next: func(c *fiber.Ctx) bool {
			return true
		},
	}))

	middleware.Validator = validator.New()

}

func setHealthCheckRoutes(app *fiber.App) {
	app.Get("/200", func(ctx *fiber.Ctx) error {
		return ctx.JSON("200")
	})
}

func setHandlers(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Initialize rest clients
	clients.InitAccountClient()
	clients.InitPaymentClient()
	clients.InitNotifyClient()
	clients.InitSavingsClient()

	merchantRep := merchant.NewRepo()
	merchantSrv := merchant.NewService(merchantRep)

	locationRep := location.NewRepo()
	locationSrv := location.NewService(locationRep)

	paymentRep := payment.NewRepo()
	//paymentSrv := payment.NewService(paymentRep)

	earningRep := earning.NewRepo()
	earningSrv := earning.NewService(earningRep)

	earningAccTxRep := earning_account_transaction.NewRepo()

	earningAccRep := earning_account.NewRepo()
	earningAccSrv := earning_account.NewService(earningAccRep, earningAccTxRep)

	mpesaStoreRep := mpesa_store.NewRepo()
	mpesaStoreSrv := mpesa_store.NewService(mpesaStoreRep)

	transactionRep := transaction.NewRepo()
	transactionSrv := transaction.NewService(transactionRep, merchantRep, paymentRep, earningAccRep, earningRep, mpesaStoreRep, earningAccSrv)

	ipnSrv := ipn.NewService(paymentRep, transactionRep, merchantRep, mpesaStoreRep, earningAccRep, earningRep, transactionSrv, earningAccSrv, earningSrv)
	jobsSrv := jobs.NewService(earningSrv)

	routes.IpnRouter(v1, ipnSrv)
	routes.JobsRouter(v1, jobsSrv)

	app.Use(jwt.New(jwt.Config{
		Secret: viper.GetString("JWT_KEY"),
		Expiry: time.Duration(15) * time.Minute,
	}))

	routes.MerchantRouter(v1, merchantSrv)
	routes.LocationRouter(v1, locationSrv)
	routes.TransactionRouter(v1, transactionSrv)
	routes.MpesaStoreRouter(v1, mpesaStoreSrv)
	routes.EarningAccountRouter(v1, earningAccSrv)
}

func Server() *fiber.App {
	// Create a new fiber instance with custom config
	app := fiber.New(fiber.Config{
		Prefork: viper.GetBool("PREFORK"),

		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.ClientLog.Error(err.Error())
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})

	// ...

	setMiddleware(app)
	setHealthCheckRoutes(app)
	setHandlers(app)

	//data, _ := json.MarshalIndent(app.GetRoutes(true), "", "  ")
	//fmt.Print(string(data))

	return app
}
