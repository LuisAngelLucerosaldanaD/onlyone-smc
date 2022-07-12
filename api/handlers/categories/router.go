package categories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/middleware"
)

func RouterCategories(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerCategories{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/categories")
	user.Get("/all", h.GetAllCategories)
	user.Post("/create-style", middleware.JWTProtected(), h.CreateStyleOfCredential)
}
