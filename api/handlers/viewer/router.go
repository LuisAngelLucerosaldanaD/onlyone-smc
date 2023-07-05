package viewer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/middleware"
)

func RouterViewer(app *fiber.App, db *sqlx.DB, txId string) {
	h := handlerViewer{Db: db, TxId: txId}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	viewer := v1.Group("viewer")
	viewer.Post("/", middleware.JWTProtected(), h.CreateShortLink)
}
