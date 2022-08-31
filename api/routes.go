package api

import (
	"github.com/ansrivas/fiberprometheus/v2"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"onlyone_smc/api/handlers/categories"
	"onlyone_smc/api/handlers/credentials"
	"onlyone_smc/api/handlers/login"
	users "onlyone_smc/api/handlers/user"
	_ "onlyone_smc/docs"
)

func routes(db *sqlx.DB, loggerHttp bool, allowedOrigins string) *fiber.App {
	app := fiber.New()

	prometheus := fiberprometheus.New("OnlyOne Smart Contract")
	prometheus.RegisterAt(app, "/metrics")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	app.Use(recover.New())
	app.Use(prometheus.Middleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST",
	}))
	if loggerHttp {
		app.Use(logger.New())
	}
	TxID := uuid.New().String()

	loadRoutes(app, db, TxID)

	return app
}

func loadRoutes(app *fiber.App, db *sqlx.DB, TxID string) {
	users.RouterCreateUser(app, db, TxID)
	login.RouterLogin(app, db, TxID)
	credentials.RouterCredentials(app, db, TxID)
	categories.RouterCategories(app, db, TxID)
}
