package viewer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/msg"
	"onlyone_smc/pkg/cfg"
)

type handlerViewer struct {
	Db   *sqlx.DB
	TxId string
}

func (h *handlerViewer) CreateShortLink(c *fiber.Ctx) error {
	res := ResponseViewer{Error: true}
	req := RequestViewer{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("No se pudo parsear el cuerpo de la petición, err: ", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.Db, h.TxId)
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	srvCfg := cfg.NewServerCfg(h.Db, nil, h.TxId)
	resPage, code, err := srvCfg.SrvCredentialPage.CreateCredentialPage(uuid.New().String(), req.Url, req.Ttl)
	if err != nil {
		logger.Error.Printf("No se pudo crear el short link, err: ", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.Db, h.TxId)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resPage == nil {
		logger.Error.Printf("No se pudo crear el short link")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.Db, h.TxId)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resPage.ID
	res.Error = false
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.Db, h.TxId)
	return c.Status(http.StatusOK).JSON(res)
}
