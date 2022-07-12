package credentials

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/middleware"
)

func RouterCredentials(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerCredentials{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/credentials")
	user.Post("/create", middleware.JWTProtected(), h.createCredential)
	user.Post("/jwt", middleware.JWTProtected(), h.getJWTTransaction)
	user.Get("/all/:block_id/:limit/:offset", middleware.JWTProtected(), h.getAllCredentials)
	user.Get("/files/:trx", middleware.JWTProtected(), h.getAllTransactionFiles)
}
