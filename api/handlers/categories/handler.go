package categories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/msg"
	"onlyone_smc/pkg/auth"
	"onlyone_smc/pkg/cfg"
	"strings"
)

type handlerCategories struct {
	DB   *sqlx.DB
	TxID string
}

func (h *handlerCategories) GetAllCategories(c *fiber.Ctx) error {
	res := responseAllCategories{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	categories, err := srvCfg.SrvCategories.GetAllCategories()
	if err != nil {
		logger.Error.Printf("couldn't get credentials: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(15, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	for _, cat := range categories {
		category := resCategories{
			Category: cat,
		}
		styles, code, err := srvCfg.SrvStyles.GetCredentialStylesByCredentialID(cat.ID)
		if err != nil {
			logger.Error.Printf("couldn't get styles of credentials: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		for _, style := range styles {

			backgroundImg := ""
			logoImg := ""

			Background := strings.Split(style.Background, "sep/")
			backgroundFile, _, err := srvAuth.SrvFiles.GetFileByPath(Background[0], Background[1])
			if err != nil {
				logger.Error.Printf("couldn't get background picture: %v", err)
				backgroundImg = ""
			}

			if backgroundFile != nil {
				backgroundImg = backgroundFile.Encoding
			}

			Logo := strings.Split(style.Logo, "sep/")
			logoFile, _, err := srvAuth.SrvFiles.GetFileByPath(Logo[0], Logo[1])
			if err != nil {
				logger.Error.Printf("couldn't get background picture: %v", err)
				logoImg = ""
			}

			if logoFile != nil {
				logoImg = logoFile.Encoding
			}

			category.Styles = append(category.Styles, Styles{
				Type:       style.Type,
				Background: backgroundImg,
				Logo:       logoImg,
				Identifiers: []identifiers{
					{
						Type:       "front",
						Attributes: style.Front,
					},
					{
						Type:       "back",
						Attributes: style.Back,
					},
				},
			})
		}

		res.Data = append(res.Data, &category)
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerCategories) CreateStyleOfCredential(c *fiber.Ctx) error {
	res := resAny{Error: true}
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	m := requestCreateStyle{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model requestCreateStyle: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileBackground, err := srvAuth.SrvFiles.UploadFile(2011, "background.jpg", m.Background)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	fileLogo, err := srvAuth.SrvFiles.UploadFile(2011, "logo.jpg", m.Logo)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err := srvCfg.SrvStyles.CreateCredentialStyles(uuid.New().String(), m.Type, fileBackground.Path+"sep/"+fileBackground.FileName, fileLogo.Path+"sep/"+fileLogo.FileName, m.Front, m.Back, m.CategoryID)
	if err != nil {
		logger.Error.Printf("couldn't create category style: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "create category"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
