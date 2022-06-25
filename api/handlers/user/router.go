package users

import (
	"onlyone_smc/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterCreateUser(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerUser{DB: db, TxID: txID}
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/user")
	user.Post("/create", h.createUser)
	user.Post("/update-password", middleware.JWTProtected(), h.changePassword)
	user.Get("/validate-email/:email", h.validateEmail)
	user.Get("/validate-nickname/:nickname", h.validateNickname)
	user.Get("/validate-identity-number/:inumber", h.getUserByIdentityNumber)
	user.Post("/validate-identity", middleware.JWTProtected(), h.validateIdentity)
	user.Get("/picture-profile", middleware.JWTProtected(), h.getUserPictureProfile)
	user.Get("/:id", h.getUserById)
}
