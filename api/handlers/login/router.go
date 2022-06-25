package login

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterLogin(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerLogin{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/login", h.Login)
}
