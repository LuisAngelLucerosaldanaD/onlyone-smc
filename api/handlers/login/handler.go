package login

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"onlyone_smc/internal/env"
	"onlyone_smc/internal/grpc/auth_proto"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/msg"
)

type handlerLogin struct {
	DB   *sqlx.DB
	TxID string
}

// Login godoc
// @Summary Login OnlyOne Smart Contract
// @Description Login OnlyOne Smart Contract
// @tags Authentication
// @Accept  json
// @Produce  json
// @Param Login body RequestLogin true "Request login"
// @Success 200 {object} ResponseLogin
// @Router /api/v1/login [post]
func (h *handlerLogin) Login(c *fiber.Ctx) error {

	res := ResponseLogin{Error: true}
	m := RequestLogin{}
	e := env.NewConfiguration()
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientAuth := auth_proto.NewAuthServicesUsersClient(connAuth)

	resLogin, err := clientAuth.Login(context.Background(), &auth_proto.LoginRequest{
		Email:    &m.Email,
		Nickname: &m.Nickname,
		Password: m.Password,
	})
	if err != nil {
		logger.Error.Printf("No se pudo obtener el token de autenticacion: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resLogin == nil {
		logger.Error.Printf("No se pudo obtener el token de autenticacion")
		res.Code, res.Type, res.Msg = msg.GetByCode(10, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resLogin.Error {
		logger.Error.Printf(resLogin.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(10, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	authRes := Token{AccessToken: resLogin.Data.AccessToken, RefreshToken: resLogin.Data.RefreshToken}
	res.Data = authRes
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
