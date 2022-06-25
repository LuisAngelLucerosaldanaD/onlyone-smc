package helpers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"onlyone_smc/internal/models"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func GetUserContext(c *fiber.Ctx) *models.User {
	userUrl := c.Locals("user").(*jwt.Token)
	claims := userUrl.Claims.(jwt.MapClaims)
	for i, cl := range claims {
		if i == "user" {
			u := models.User{}
			ub, _ := json.Marshal(cl)
			_ = json.Unmarshal(ub, &u)
			return &u
		}
	}
	return nil
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}
	return buff.Bytes()
}
