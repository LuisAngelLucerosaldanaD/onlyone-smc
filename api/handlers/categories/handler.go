package categories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/msg"
	"onlyone_smc/pkg/cfg"
)

type handlerCategories struct {
	DB   *sqlx.DB
	TxID string
}

func (h *handlerCategories) GetAllCategories(c *fiber.Ctx) error {
	res := responseAllCategories{Error: true}

	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	categories, err := srvCfg.SrvCategories.GetAllCategories()
	if err != nil {
		logger.Error.Printf("couldn't get transactions: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(15, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = categories
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
