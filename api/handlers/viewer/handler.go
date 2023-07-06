package viewer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
	"onlyone_smc/internal/logger"
	"onlyone_smc/pkg/cfg"
	"time"
)

type handlerViewer struct {
	Db   *sqlx.DB
	TxId string
}

func (h *handlerViewer) GetShortLink(c *fiber.Ctx) error {
	idCtx := c.Params("id")
	if idCtx == "" {
		return c.Status(http.StatusNotFound).SendString("No se encontro la pagina solicitada")
	}
	srvCfg := cfg.NewServerCfg(h.Db, nil, h.TxId)
	resPage, _, err := srvCfg.SrvCredentialPage.GetCredentialPageByID(idCtx)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el short link, err: ", err.Error())
		return c.Status(http.StatusAccepted).SendString("No se pudo obtener la pagina solicitada")
	}

	if resPage == nil {
		logger.Error.Printf("No se pudo obtener el Short Link")
		return c.Status(http.StatusNotFound).SendString("No se encontro la pagina solicitada")
	}

	lifeTtl := time.Now().Sub(time.Now()).Seconds()
	if int(lifeTtl) < (resPage.Ttl * 1000) {
		_, err = srvCfg.SrvCredentialPage.DeleteCredentialPage(resPage.ID)
		if err != nil {
			logger.Error.Printf("No se pudo obtener el short link, err: ", err.Error())
			return c.Status(http.StatusAccepted).SendString("No se pudo obtener la pagina")
		}
		return c.Status(http.StatusAccepted).SendString("La pagina solicita ha caducado o ya no se encuentra disponible")
	}

	return c.Redirect(resPage.Url, fiber.StatusMovedPermanently)
}
